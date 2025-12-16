package operations

import (
	"fmt"
	"time"

	"github.com/gooddata/gooddata-neobackstop/internals"
	"github.com/playwright-community/playwright-go"
)

func ScrollToSelector(logPrefix string, page playwright.Page, scenario internals.Scenario) *string {
	if scenario.ScrollToSelector != nil {
		sts := *scenario.ScrollToSelector
		t0 := time.Now()
		fmt.Println(logPrefix, "scrollToSelector: waiting for", sts)

		// Wait for the element first
		_, err := page.WaitForSelector(sts, playwright.PageWaitForSelectorOptions{
			State:   playwright.WaitForSelectorStateVisible,
			Timeout: playwright.Float(10000),
		})
		if err != nil {
			e := "ScrollToSelector " + sts + " not found"
			return &e
		}
		waitTime := time.Since(t0).Milliseconds()

		tScroll := time.Now()
		// Scroll element into view
		_, err = page.EvalOnSelector(sts, `el => el.scrollIntoView({ behavior: "instant", block: "center" })`, nil)
		if err != nil {
			e := "Failed to scroll element into view: " + err.Error()
			return &e
		}
		scrollTime := time.Since(tScroll).Milliseconds()
		fmt.Println(logPrefix, "scrollToSelector:", sts, "completed in", waitTime+scrollTime, "ms (wait:", waitTime, "ms, scroll:", scrollTime, "ms)")

		return postInteractionWait(logPrefix, page, scenario.PostInteractionWait)
	}

	return nil
}
