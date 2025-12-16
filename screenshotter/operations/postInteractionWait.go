package operations

import (
	"fmt"
	"time"

	"github.com/gooddata/gooddata-neobackstop/scenario"
	"github.com/playwright-community/playwright-go"
)

func postInteractionWait(logPrefix string, page playwright.Page, piw *scenario.SelectorThenDelay) *string {
	if piw != nil {
		if piw.Selector != nil {
			// selector, wait for it
			selector := *piw.Selector
			t0 := time.Now()

			_, err := page.WaitForSelector(selector, playwright.PageWaitForSelectorOptions{
				State: playwright.WaitForSelectorStateAttached,
				// use default 30s timeout
			})
			if err != nil {
				e := "PostInteractionWait " + selector + " didn't appear"
				return &e
			}
			fmt.Println(logPrefix, "postInteractionWait selector", selector, "appeared in", time.Since(t0).Milliseconds(), "ms")
		}
		if piw.Delay != nil {
			// delay, wait
			delayMs := piw.Delay.Milliseconds()
			fmt.Println(logPrefix, "postInteractionDelay: starting sleep for", delayMs, "ms")
			time.Sleep(*piw.Delay)
			fmt.Println(logPrefix, "postInteractionDelay: ending sleep for", delayMs, "ms")
		}
	}

	return nil
}
