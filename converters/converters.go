package converters

import (
	"fmt"
	"sort"
	"time"

	"github.com/gooddata/gooddata-neobackstop/browser"
	"github.com/gooddata/gooddata-neobackstop/internals"
	"github.com/gooddata/gooddata-neobackstop/scenario"
	"github.com/gooddata/gooddata-neobackstop/viewport"
)

func convertSelectorWithBeforeAfterDelay(selectors []interface{}) []internals.SelectorWithBeforeAfterDelay {
	if selectors == nil {
		return nil
	}

	issTmp := make([]*internals.SelectorWithBeforeAfterDelay, 0) // slice of pointers to make manipulation easier

	var firstTimeout *time.Duration
	for _, selector := range selectors {
		switch s := selector.(type) {
		case string:
			// legacy format: selector
			is := internals.SelectorWithBeforeAfterDelay{Selector: s}
			if firstTimeout != nil {
				is.WaitBefore = firstTimeout
				firstTimeout = nil
			}
			issTmp = append(issTmp, &is)
		case float64: // json number fields unmarshalled into interfaces are float64s
			// legacy format: timeout
			if len(issTmp) == 0 {
				d := time.Duration(s) * time.Millisecond
				if firstTimeout != nil {
					// starts with multiple timeouts...
					// corner case but we will handle it anyway
					d += *firstTimeout
				}
				firstTimeout = &d
			} else {
				lastIndex := len(issTmp) - 1
				lastIs := issTmp[lastIndex]

				d := time.Duration(s) * time.Millisecond
				if lastIs.WaitAfter != nil {
					// two consecutive delays
					// another conner case that we will handle
					d += *lastIs.WaitAfter
				}
				lastIs.WaitAfter = &d
			}
		case map[string]interface{}:
			// new format: object with selector, waitBefore, waitAfter
			is := internals.SelectorWithBeforeAfterDelay{}
			if v, ok := s["selector"].(string); ok {
				is.Selector = v
			}
			if v, ok := s["waitBefore"].(float64); ok {
				d := time.Duration(v) * time.Millisecond
				is.WaitBefore = &d
			}
			if v, ok := s["waitAfter"].(float64); ok {
				d := time.Duration(v) * time.Millisecond
				is.WaitAfter = &d
			}
			issTmp = append(issTmp, &is)
		default:
			fmt.Println(selectors)
			panic("Unknown click/hover selector type")
		}
	}

	// convert slice of pointers to regular slice
	iss := make([]internals.SelectorWithBeforeAfterDelay, len(issTmp))
	for i, ih := range issTmp {
		iss[i] = *ih
	}
	return iss
}

func convertSelectorOrDelay(value interface{}) *internals.SelectorThenDelay {
	if value == nil {
		return nil
	}

	var sod internals.SelectorThenDelay
	// check type
	switch piw := value.(type) {
	case string:
		// legacy format: selector
		sod = internals.SelectorThenDelay{
			Selector: &piw,
		}
	case float64:
		// legacy format: timeout
		d := time.Duration(piw) * time.Millisecond
		sod = internals.SelectorThenDelay{
			Delay: &d,
		}
	case map[string]interface{}:
		// new format: object with selector, delay
		if v, ok := piw["selector"].(string); ok {
			sod.Selector = &v
		}
		if v, ok := piw["delay"].(float64); ok {
			d := time.Duration(v) * time.Millisecond
			sod.Delay = &d
		}
	default:
		fmt.Println(piw)
		panic("Unknown PostInteractionWait type")
	}

	return &sod
}

func scenarioToInternal(b browser.Browser, v viewport.Viewport, s scenario.Scenario) internals.Scenario {
	return internals.Scenario{
		Browser:             b,
		Viewport:            v,
		Id:                  s.Id,
		Label:               s.Label,
		Url:                 s.Url,
		ReadySelector:       s.ReadySelector,
		Delay:               s.Delay,
		ReloadAfterReady:    s.ReloadAfterReady,
		KeyPressSelector:    (*internals.KeyPressSelector)(s.KeyPressSelector),
		HoverSelector:       s.HoverSelector,
		HoverSelectors:      convertSelectorWithBeforeAfterDelay(s.HoverSelectors),
		ClickSelector:       s.ClickSelector,
		ClickSelectors:      convertSelectorWithBeforeAfterDelay(s.ClickSelectors),
		PostInteractionWait: convertSelectorOrDelay(s.PostInteractionWait),
		ScrollToSelector:    s.ScrollToSelector,
		MisMatchThreshold:   s.MisMatchThreshold,
	}
}

func ScenariosToInternal(browsers []browser.Browser, viewports []viewport.Viewport, scenarios []scenario.Scenario) []internals.Scenario {
	output := make([]internals.Scenario, 0) // we could pre-calculate this, but until we do multi-browser testing, it's not worth it

	for _, s := range scenarios {
		if s.Browsers == nil {
			for _, b := range browsers {
				if s.Viewports == nil {
					for _, v := range viewports {
						output = append(output, scenarioToInternal(b, v, s))
					}
				} else {
					for _, v := range s.Viewports {
						output = append(output, scenarioToInternal(b, v, s))
					}
				}
			}
		} else {
			for _, b := range s.Browsers {
				if s.Viewports == nil {
					for _, v := range viewports {
						output = append(output, scenarioToInternal(b, v, s))
					}
				} else {
					for _, v := range s.Viewports {
						output = append(output, scenarioToInternal(b, v, s))
					}
				}
			}
		}
	}

	// sort the internalScenarios by browser (a to z), then by viewport (width first, smallest to largest)
	sort.Slice(output, func(i, j int) bool {
		// 1. Sort by Browser (alphabetically)
		if output[i].Browser != output[j].Browser {
			return output[i].Browser < output[j].Browser
		}
		// 2. Sort by Viewport.Width (smallest first)
		if output[i].Viewport.Width != output[j].Viewport.Width {
			return output[i].Viewport.Width < output[j].Viewport.Width
		}
		// 3. Sort by Viewport.Height (smallest first)
		return output[i].Viewport.Height < output[j].Viewport.Height
	})

	return output
}
