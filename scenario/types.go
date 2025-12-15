package scenario

import (
	"encoding/json"
	"time"

	"github.com/gooddata/gooddata-neobackstop/browser"
	"github.com/gooddata/gooddata-neobackstop/viewport"
)

type Delay struct {
	PostReady     time.Duration `json:"postReady"`
	PostOperation time.Duration `json:"postOperation"`
}

func (d *Delay) UnmarshalJSON(data []byte) error {
	var raw struct {
		PostReady     float64 `json:"postReady"`
		PostOperation float64 `json:"postOperation"`
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	d.PostReady = time.Duration(raw.PostReady) * time.Millisecond
	d.PostOperation = time.Duration(raw.PostOperation) * time.Millisecond
	return nil
}

type KeyPressSelector struct {
	KeyPress string `json:"keyPress"`
	Selector string `json:"selector"`
}

type SelectorWithBeforeAfterDelay struct {
	Selector   string         `json:"selector"`
	WaitBefore *time.Duration `json:"waitBefore"`
	WaitAfter  *time.Duration `json:"waitAfter"`
}

func (d *SelectorWithBeforeAfterDelay) UnmarshalJSON(data []byte) error {
	var raw struct {
		Selector   string   `json:"selector"`
		WaitBefore *float64 `json:"waitBefore"`
		WaitAfter  *float64 `json:"waitAfter"`
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	d.Selector = raw.Selector
	if raw.WaitBefore != nil {
		t := time.Duration(*raw.WaitBefore) * time.Millisecond
		d.WaitBefore = &t
	}
	if raw.WaitAfter != nil {
		t := time.Duration(*raw.WaitAfter) * time.Millisecond
		d.WaitAfter = &t
	}
	return nil
}

// Scenario - properties in order of processing, scenarios.json must be an array of Scenario
type Scenario struct {
	Browsers            []browser.Browser              `json:"browsers"`
	Viewports           []viewport.Viewport            `json:"viewports"`
	Id                  string                         `json:"id"`
	Label               string                         `json:"label"`
	Url                 string                         `json:"url"`
	ReadySelector       *string                        `json:"readySelector"`
	ReloadAfterReady    bool                           `json:"reloadAfterReady"`
	Delay               *Delay                         `json:"delay"`
	KeyPressSelector    *KeyPressSelector              `json:"keyPressSelector"`
	HoverSelector       *string                        `json:"hoverSelector"`
	HoverSelectors      []SelectorWithBeforeAfterDelay `json:"hoverSelectors"`
	ClickSelector       *string                        `json:"clickSelector"`
	ClickSelectors      []SelectorWithBeforeAfterDelay `json:"clickSelectors"`
	PostInteractionWait interface{}                    `json:"postInteractionWait"`
	ScrollToSelector    *string                        `json:"scrollToSelector"`
	MisMatchThreshold   *float64                       `json:"misMatchThreshold"`
}
