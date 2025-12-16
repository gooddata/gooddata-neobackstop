package operations

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gooddata/gooddata-neobackstop/internals"
	"github.com/playwright-community/playwright-go"
)

func ClickSelectors(logPrefix string, page playwright.Page, scenario internals.Scenario) *string {
	if scenario.ClickSelectors != nil {
		fmt.Println(logPrefix, "clickSelectors: processing", len(scenario.ClickSelectors), "clicks")

		for i, cs := range scenario.ClickSelectors {
			idx := strconv.Itoa(i + 1)

			if cs.WaitBefore != nil {
				delayMs := cs.WaitBefore.Milliseconds()
				fmt.Println(logPrefix, "clickSelectors["+idx+"]: waitBefore starting sleep for", delayMs, "ms")
				time.Sleep(*cs.WaitBefore)
				fmt.Println(logPrefix, "clickSelectors["+idx+"]: waitBefore ending sleep for", delayMs, "ms")
			}

			t0 := time.Now()
			fmt.Println(logPrefix, "clickSelectors["+idx+"]: waiting for", cs.Selector)

			_, err := page.WaitForSelector(cs.Selector, playwright.PageWaitForSelectorOptions{
				State:   playwright.WaitForSelectorStateVisible,
				Timeout: playwright.Float(10000), // default 10s
			})
			if err != nil {
				e := "ClickSelector " + cs.Selector + " didn't appear"
				return &e
			}
			waitTime := time.Since(t0).Milliseconds()

			tClick := time.Now()
			err = page.Click(cs.Selector)
			if err != nil {
				e := "ClickSelector " + cs.Selector + " couldn't be clicked"
				return &e
			}
			clickTime := time.Since(tClick).Milliseconds()
			fmt.Println(logPrefix, "clickSelectors["+idx+"]:", cs.Selector, "completed in", waitTime+clickTime, "ms (wait:", waitTime, "ms, click:", clickTime, "ms)")

			sErr := postInteractionWait(logPrefix, page, scenario.PostInteractionWait)
			if sErr != nil {
				return sErr
			}

			if cs.WaitAfter != nil {
				delayMs := cs.WaitAfter.Milliseconds()
				fmt.Println(logPrefix, "clickSelectors["+idx+"]: waitAfter starting sleep for", delayMs, "ms")
				time.Sleep(*cs.WaitAfter)
				fmt.Println(logPrefix, "clickSelectors["+idx+"]: waitAfter ending sleep for", delayMs, "ms")
			}
		}
	}

	return nil
}
