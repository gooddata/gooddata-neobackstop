package internals

import (
	"time"

	"github.com/gooddata/gooddata-neobackstop/browser"
	"github.com/gooddata/gooddata-neobackstop/scenario"
	"github.com/gooddata/gooddata-neobackstop/viewport"
)

type SelectorThenDelay struct {
	Selector *string        `json:"selector"`
	Delay    *time.Duration `json:"delay"`
}

// Scenario - An internal scenario type, properties in order of processing, constructed from scenario.Scenario and config.Config
// this type is exposed in ci-report so needs json keys
type Scenario struct {
	Browser             browser.Browser                         `json:"browser"`
	Viewport            viewport.Viewport                       `json:"viewport"`
	Id                  string                                  `json:"id"`
	Label               string                                  `json:"label"`
	Url                 string                                  `json:"url"`
	ReadySelector       *string                                 `json:"readySelector"`
	ReloadAfterReady    bool                                    `json:"reloadAfterReady"`
	Delay               *scenario.Delay                         `json:"delay"`
	KeyPressSelector    *scenario.KeyPressSelector              `json:"keyPressSelector"`
	HoverSelector       *string                                 `json:"hoverSelector"`
	HoverSelectors      []scenario.SelectorWithBeforeAfterDelay `json:"hoverSelectors"`
	ClickSelector       *string                                 `json:"clickSelector"`
	ClickSelectors      []scenario.SelectorWithBeforeAfterDelay `json:"clickSelectors"`
	PostInteractionWait *SelectorThenDelay                      `json:"postInteractionWait"`
	ScrollToSelector    *string                                 `json:"scrollToSelector"`
	MisMatchThreshold   *float64                                `json:"misMatchThreshold"`
}
