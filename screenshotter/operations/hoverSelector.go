package operations

import (
	"fmt"
	"time"

	"github.com/gooddata/gooddata-neobackstop/internals"
	"github.com/playwright-community/playwright-go"
)

func HoverSelector(logPrefix string, page playwright.Page, scenario internals.Scenario) *string {
	if scenario.HoverSelector != nil {
		hs := *scenario.HoverSelector
		t0 := time.Now()
		fmt.Println(logPrefix, "hoverSelector: waiting for", hs)

		// there is a HoverSelector, wait for it
		_, err := page.WaitForSelector(hs, playwright.PageWaitForSelectorOptions{
			State: playwright.WaitForSelectorStateAttached,
			// use default 30s timeout
		})
		if err != nil {
			e := "HoverSelector " + hs + " didn't appear"
			return &e
		}
		waitTime := time.Since(t0).Milliseconds()

		tHover := time.Now()
		err = page.Hover(hs)
		if err != nil {
			e := "HoverSelector " + hs + "couldn't be hovered"
			return &e
		}
		hoverTime := time.Since(tHover).Milliseconds()
		fmt.Println(logPrefix, "hoverSelector:", hs, "completed in", waitTime+hoverTime, "ms (wait:", waitTime, "ms, hover:", hoverTime, "ms)")

		return postInteractionWait(logPrefix, page, scenario.PostInteractionWait)
	}

	return nil
}
