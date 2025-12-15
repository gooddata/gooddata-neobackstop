package scenario

import (
	"github.com/gooddata/gooddata-neobackstop/browser"
	"github.com/gooddata/gooddata-neobackstop/internals"
	"github.com/gooddata/gooddata-neobackstop/viewport"
)

// Scenario - properties in order of processing, scenarios.json must be an array of Scenario
type Scenario struct {
	Browsers            []browser.Browser           `json:"browsers"`
	Viewports           []viewport.Viewport         `json:"viewports"`
	Id                  string                      `json:"id"`
	Label               string                      `json:"label"`
	Url                 string                      `json:"url"`
	ReadySelector       *string                     `json:"readySelector"`
	ReloadAfterReady    bool                        `json:"reloadAfterReady"`
	Delay               *internals.Delay            `json:"delay"`
	KeyPressSelector    *internals.KeyPressSelector `json:"keyPressSelector"`
	HoverSelector       *string                     `json:"hoverSelector"`
	HoverSelectors      []interface{}               `json:"hoverSelectors"`
	ClickSelector       *string                     `json:"clickSelector"`
	ClickSelectors      []interface{}               `json:"clickSelectors"`
	PostInteractionWait interface{}                 `json:"postInteractionWait"`
	ScrollToSelector    *string                     `json:"scrollToSelector"`
	MisMatchThreshold   *float64                    `json:"misMatchThreshold"`
}
