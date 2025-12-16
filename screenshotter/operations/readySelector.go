package operations

import (
	"fmt"
	"time"

	"github.com/gooddata/gooddata-neobackstop/scenario"
	"github.com/playwright-community/playwright-go"
)

func ReadySelector(logPrefix string, page playwright.Page, value *scenario.ReadySelector) *string {
	if value != nil {
		state := playwright.WaitForSelectorState(value.State)
		t0 := time.Now()
		fmt.Println(logPrefix, "readySelector: waiting for", value.Selector, "state:", value.State)

		// there is a ReadySelector, wait for it
		_, err := page.WaitForSelector(value.Selector, playwright.PageWaitForSelectorOptions{
			State:   &state,
			Timeout: playwright.Float(30000), // todo: add timeout to config
		})
		if err != nil {
			e := "ReadySelector " + value.Selector + " didn't appear"
			return &e
		}
		fmt.Println(logPrefix, "readySelector:", value.Selector, "appeared in", time.Since(t0).Milliseconds(), "ms")
	}

	return nil
}
