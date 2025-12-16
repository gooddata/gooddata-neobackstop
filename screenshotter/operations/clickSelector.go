package operations

import (
	"fmt"
	"time"

	"github.com/gooddata/gooddata-neobackstop/internals"
	"github.com/playwright-community/playwright-go"
)

func ClickSelector(logPrefix string, page playwright.Page, scenario internals.Scenario) *string {
	if scenario.ClickSelector != nil {
		cs := *scenario.ClickSelector
		t0 := time.Now()
		fmt.Println(logPrefix, "clickSelector: waiting for", cs)

		_, err := page.WaitForSelector(cs, playwright.PageWaitForSelectorOptions{
			State:   playwright.WaitForSelectorStateVisible,
			Timeout: playwright.Float(10000), // default 10s
		})
		if err != nil {
			e := "ClickSelector " + cs + " didn't appear"
			return &e
		}
		waitTime := time.Since(t0).Milliseconds()

		tClick := time.Now()
		err = page.Click(cs)
		if err != nil {
			e := "ClickSelector " + cs + " couldn't be clicked"
			return &e
		}
		clickTime := time.Since(tClick).Milliseconds()
		fmt.Println(logPrefix, "clickSelector:", cs, "completed in", waitTime+clickTime, "ms (wait:", waitTime, "ms, click:", clickTime, "ms)")

		return postInteractionWait(logPrefix, page, scenario.PostInteractionWait)
	}

	return nil
}
