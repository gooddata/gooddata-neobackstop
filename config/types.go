package config

import (
	"github.com/gooddata/gooddata-neobackstop/browser"
	"github.com/gooddata/gooddata-neobackstop/viewport"
)

type HtmlReportConfig struct {
	Path                string `json:"path"`
	ShowSuccessfulTests bool   `json:"showSuccessfulTests"`
}

type BrowserSettings struct {
	Name browser.Browser `json:"name"`
	Args []string        `json:"args"`
}

type Config struct {
	Browsers              map[string]BrowserSettings `json:"browsers"`
	DefaultBrowsers       []string                   `json:"defaultBrowsers"`
	Viewports             []viewport.Viewport        `json:"viewports"`
	BitmapsReferencePath  string                     `json:"bitmapsReferencePath"`
	BitmapsTestPath       string                     `json:"bitmapsTestPath"`
	HtmlReport            HtmlReportConfig           `json:"htmlReport"`
	CiReportPath          string                     `json:"ciReportPath"`
	AsyncCaptureLimit     int                        `json:"asyncCaptureLimit"`
	AsyncCompareLimit     int                        `json:"asyncCompareLimit"`
	RequireSameDimensions bool                       `json:"requireSameDimensions"`
	RetryCount            int                        `json:"retryCount"`
}
