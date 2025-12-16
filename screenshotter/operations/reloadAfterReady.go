package operations

import (
	"fmt"
	"log"
	"time"

	"github.com/gooddata/gooddata-neobackstop/internals"
	"github.com/playwright-community/playwright-go"
)

func ReloadAfterReady(logPrefix string, page playwright.Page, scenario internals.Scenario) *string {
	// todo: we might be able to check if the scenario runs on chromium
	//  before doing this, as these issues do not happen on firefox

	// todo: consider enabling this by default at the end, depends on how much it affects time
	t0 := time.Now()
	fmt.Println(logPrefix, "reloadAfterReady: reloading page")

	if _, err := page.Reload(playwright.PageReloadOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
	}); err != nil {
		log.Panicf("could not reload: %v", err)
	}
	fmt.Println(logPrefix, "reloadAfterReady: page reloaded in", time.Since(t0).Milliseconds(), "ms")

	return ReadySelector(logPrefix, page, scenario.ReadySelector)
}
