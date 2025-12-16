package comparer

import (
	"github.com/gooddata/gooddata-neobackstop/screenshotter"
	"github.com/gooddata/gooddata-neobackstop/utils"
)

// Result - todo some of these should be pointers
type Result struct {
	ScreenshotterResult *screenshotter.Result
	HasReference        bool
	DiffCreated         bool
	MismatchPercentage  *float64
	MatchesReference    bool
	DimensionMismatch   *utils.DimensionDiff
	Error               *string
}
