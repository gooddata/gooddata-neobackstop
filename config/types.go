package config

import (
	"github.com/gooddata/gdc-neobackstop/browser"
	"github.com/gooddata/gdc-neobackstop/viewport"
)

type HtmlReportConfig struct {
	Path                string `json:"path"`
	ShowSuccessfulTests bool   `json:"showSuccessfulTests"`
}

type Config struct {
	Browsers             []browser.Browser   `json:"browsers"`
	Viewports            []viewport.Viewport `json:"viewports"`
	BitmapsReferencePath string              `json:"bitmapsReferencePath"`
	BitmapsTestPath      string              `json:"bitmapsTestPath"`
	HtmlReport           HtmlReportConfig    `json:"htmlReport"`
	CiReportPath         string              `json:"ciReportPath"`
	Args                 []string            `json:"args"`
	AsyncCaptureLimit    int                 `json:"asyncCaptureLimit"`
	AsyncCompareLimit    int                 `json:"asyncCompareLimit"`
	// RequireSameDimensions bool                `json:"requireSameDimensions"`
}
