package converters

import (
	"sort"

	"github.com/gooddata/gooddata-neobackstop/browser"
	"github.com/gooddata/gooddata-neobackstop/internals"
	"github.com/gooddata/gooddata-neobackstop/scenario"
	"github.com/gooddata/gooddata-neobackstop/viewport"
)

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
		KeyPressSelector:    s.KeyPressSelector,
		HoverSelector:       s.HoverSelector,
		HoverSelectors:      s.HoverSelectors,
		ClickSelector:       s.ClickSelector,
		ClickSelectors:      s.ClickSelectors,
		PostInteractionWait: s.PostInteractionWait,
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
