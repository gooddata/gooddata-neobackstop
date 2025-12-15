package operations

import (
	"github.com/gooddata/gooddata-neobackstop/scenario"
	"github.com/playwright-community/playwright-go"
)

func ReadySelector(page playwright.Page, value *scenario.ReadySelector) *string {
	if value != nil {
		state := playwright.WaitForSelectorState(value.State)
		// there is a ReadySelector, wait for it
		_, err := page.WaitForSelector(value.Selector, playwright.PageWaitForSelectorOptions{
			State:   &state,
			Timeout: playwright.Float(30000), // todo: add timeout to config
		})
		if err != nil {
			e := "ReadySelector " + value.Selector + " didn't appear"
			return &e
		}
	}

	return nil
}
