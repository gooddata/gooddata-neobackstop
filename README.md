# NeoBackstop

A high-performance visual regression testing tool written in Go, powered by Playwright. NeoBackstop captures screenshots of web pages and compares them against reference images to detect visual changes.

## Features

- Multi-browser support (Chromium, Firefox)
- Parallel screenshot capture and comparison
- Rich scenario configuration with selectors and interactions
- HTML report generation (BackstopJS-compatible format)
- CI-friendly JSON output
- Docker support for consistent cross-platform execution
- Library mode for custom integrations

## Installation

NeoBackstop can be used as a standalone CLI tool or as a Go library.

### Docker (Recommended)

Pull the pre-built image from Docker Hub:

```bash
docker pull gooddata/gooddata-neobackstop:latest
```

Or build locally:

```bash
docker build -t neobackstop .
```

### Go Package

Install as a Go module for library mode or standalone usage:

```bash
go get github.com/gooddata/gooddata-neobackstop@latest
```

### Building from Source

Prerequisites:
- Go 1.25.4 or later
- Playwright browsers (installed automatically on first run)

```bash
go build -o neobackstop .
```

## Usage

NeoBackstop can be used in three ways: standalone mode, Docker mode, or as a library.

### Standalone Mode

#### Test Mode

Captures screenshots and compares them against reference images:

```bash
./neobackstop test --config=./config.json --scenarios=./scenarios.json
```

#### Approve Mode

Captures screenshots and saves them as new reference images:

```bash
./neobackstop approve --config=./config.json --scenarios=./scenarios.json
```

### Docker Mode

Using the published Docker image:

```bash
# Test mode
docker run -v $(pwd)/config:/config -v $(pwd)/output:/output \
  gooddata/gooddata-neobackstop:latest test \
  --config=/config/config.json --scenarios=/config/scenarios.json

# Approve mode
docker run -v $(pwd)/config:/config -v $(pwd)/output:/output \
  gooddata/gooddata-neobackstop:latest approve \
  --config=/config/config.json --scenarios=/config/scenarios.json
```

### Library Mode

NeoBackstop can be imported as a Go library, allowing you to build custom testing workflows, add memory monitoring, or integrate with your existing test infrastructure.

#### Basic Library Usage

```go
package main

import (
    "encoding/json"
    "io"
    "log"
    "os"
    "sync"

    "github.com/gooddata/gooddata-neobackstop/comparer"
    "github.com/gooddata/gooddata-neobackstop/config"
    "github.com/gooddata/gooddata-neobackstop/converters"
    "github.com/gooddata/gooddata-neobackstop/internals"
    "github.com/gooddata/gooddata-neobackstop/scenario"
    "github.com/gooddata/gooddata-neobackstop/screenshotter"
    "github.com/playwright-community/playwright-go"
)

func main() {
    // Load configuration
    configFile, _ := os.Open("config.json")
    configBytes, _ := io.ReadAll(configFile)
    var cfg config.Config
    json.Unmarshal(configBytes, &cfg)

    // Load scenarios
    scenariosFile, _ := os.Open("scenarios.json")
    scenariosBytes, _ := io.ReadAll(scenariosFile)
    var scenarios []scenario.Scenario
    json.Unmarshal(scenariosBytes, &scenarios)

    // Convert to internal format
    internalScenarios := converters.ScenariosToInternal(
        cfg.Browsers, cfg.Viewports, scenarios,
    )

    // Install and run Playwright
    browsers := make([]string, len(cfg.Browsers))
    for i, b := range cfg.Browsers {
        browsers[i] = string(b)
    }
    playwright.Install(&playwright.RunOptions{Browsers: browsers})
    pw, _ := playwright.Run()

    // Set up worker pool for screenshots
    jobs := make(chan internals.Scenario, len(internalScenarios))
    results := make(chan screenshotter.Result, len(internalScenarios))
    var wg sync.WaitGroup

    for w := 1; w <= cfg.AsyncCaptureLimit; w++ {
        wg.Add(1)
        go screenshotter.Run("./output", pw, cfg, jobs, &wg, results, w)
    }

    // Send jobs
    for _, s := range internalScenarios {
        jobs <- s
    }
    close(jobs)
    wg.Wait()
    close(results)

    pw.Stop()
}
```

#### Debug Mode (Non-Headless)

For debugging scenarios locally, you can run the browser in non-headless mode to see exactly what's happening:

```go
package main

import (
    "encoding/json"
    "io"
    "os"

    "github.com/gooddata/gooddata-neobackstop/config"
    "github.com/gooddata/gooddata-neobackstop/converters"
    "github.com/gooddata/gooddata-neobackstop/internals"
    "github.com/gooddata/gooddata-neobackstop/scenario"
    "github.com/gooddata/gooddata-neobackstop/screenshotter"
    "github.com/playwright-community/playwright-go"
)

func main() {
    // Load config and scenarios (abbreviated)
    configFile, _ := os.Open("config.json")
    configBytes, _ := io.ReadAll(configFile)
    var cfg config.Config
    json.Unmarshal(configBytes, &cfg)

    scenariosFile, _ := os.Open("scenarios.json")
    scenariosBytes, _ := io.ReadAll(scenariosFile)
    var scenarios []scenario.Scenario
    json.Unmarshal(scenariosBytes, &scenarios)

    internalScenarios := converters.ScenariosToInternal(
        cfg.Browsers, cfg.Viewports, scenarios,
    )

    // Find a specific scenario to debug
    var debugScenario *internals.Scenario
    for _, s := range internalScenarios {
        if s.Label == "Dashboard/Chart View" {
            debugScenario = &s
            break
        }
    }

    // Launch browser in NON-HEADLESS mode for debugging
    pw, _ := playwright.Run()
    browser, _ := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
        Headless: playwright.Bool(false), // Show the browser!
        Args:     cfg.Args["chromium"],
    })

    context, _ := browser.NewContext(playwright.BrowserNewContextOptions{
        Viewport: &playwright.Size{
            Width:  debugScenario.Viewport.Width,
            Height: debugScenario.Viewport.Height,
        },
    })
    page, _ := context.NewPage()

    // Create a results channel (for the Job function)
    results := make(chan screenshotter.Result)
    go func() {
        for range results {} // Drain results
    }()

    // Run the screenshot job with debug mode enabled
    screenshotter.Job("./debug-output", debugScenario.Viewport.Label,
        page, *debugScenario, results, true) // true = debug mode

    browser.Close()
    pw.Stop()
}
```

#### Available Packages

When using NeoBackstop as a library, the following packages are available:

| Package         | Description                                        |
|-----------------|----------------------------------------------------|
| `config`        | Configuration types (`Config`, `HtmlReportConfig`) |
| `scenario`      | Scenario types for JSON parsing                    |
| `internals`     | Internal scenario representation after conversion  |
| `converters`    | Convert scenarios to internal format               |
| `screenshotter` | Screenshot capture worker and job functions        |
| `comparer`      | Image comparison worker                            |
| `browser`       | Browser enum (`Chromium`, `Firefox`)               |
| `viewport`      | Viewport type definition                           |
| `result`        | Result types for CI output                         |
| `html_report`   | HTML report types                                  |
| `utils`         | Utility functions                                  |

## Configuration

### config.json

The main configuration file controls browser settings, viewports, output paths, and concurrency.

```json
{
    "id": "my-visual-tests",
    "browsers": ["chromium", "firefox"],
    "viewports": [
        {
            "label": "desktop",
            "width": 1024,
            "height": 768
        },
        {
            "label": "mobile",
            "width": 375,
            "height": 667
        }
    ],
    "bitmapsReferencePath": "./output/reference",
    "bitmapsTestPath": "./output/test",
    "htmlReport": {
        "path": "./output/html-report",
        "showSuccessfulTests": false
    },
    "ciReportPath": "./output/ci-report",
    "args": {
        "chromium": [
            "--disable-infobars",
            "--disable-background-networking",
            "--disable-background-timer-throttling",
            "--disable-backgrounding-occluded-windows",
            "--disable-breakpad",
            "--disable-client-side-phishing-detection",
            "--disable-default-apps",
            "--disable-dev-shm-usage",
            "--disable-extensions",
            "--disable-features=site-per-process",
            "--disable-hang-monitor",
            "--disable-ipc-flooding-protection",
            "--disable-popup-blocking",
            "--disable-prompt-on-repost",
            "--disable-renderer-backgrounding",
            "--disable-sync",
            "--disable-translate",
            "--metrics-recording-only",
            "--no-first-run",
            "--safebrowsing-disable-auto-update",
            "--enable-automation",
            "--disable-component-update",
            "--disable-web-resource",
            "--mute-audio",
            "--no-sandbox",
            "--disable-software-rasterizer",
            "--disable-gpu",
            "--disable-setuid-sandbox",
            "--force-device-scale-factor=1"
        ],
        "firefox": [
            "--disable-dev-shm-usage",
            "--disable-extensions",
            "--enable-automation",
            "--mute-audio",
            "--no-sandbox",
            "--disable-gpu"
        ]
    },
    "asyncCaptureLimit": 2,
    "asyncCompareLimit": 6
}
```

#### Configuration Options

| Option                           | Type       | Description                                |
|----------------------------------|------------|--------------------------------------------|
| `id`                             | string     | Identifier for the test suite              |
| `browsers`                       | string[]   | Browsers to use: `"chromium"`, `"firefox"` |
| `viewports`                      | Viewport[] | List of viewport configurations            |
| `bitmapsReferencePath`           | string     | Path to store reference screenshots        |
| `bitmapsTestPath`                | string     | Path to store test screenshots             |
| `htmlReport.path`                | string     | Path for HTML report output                |
| `htmlReport.showSuccessfulTests` | boolean    | Include passing tests in HTML report       |
| `ciReportPath`                   | string     | Path for CI JSON report                    |
| `args`                           | object     | Browser-specific launch arguments          |
| `asyncCaptureLimit`              | number     | Max concurrent screenshot captures         |
| `asyncCompareLimit`              | number     | Max concurrent image comparisons           |

#### Viewport Configuration

| Property | Type   | Description                  |
|----------|--------|------------------------------|
| `label`  | string | Human-readable viewport name |
| `width`  | number | Viewport width in pixels     |
| `height` | number | Viewport height in pixels    |

### scenarios.json

Defines the test scenarios - which pages to capture and how to interact with them.

```json
[
    {
        "id": "homepage",
        "label": "Homepage/Default View",
        "url": "http://localhost:3000/",
        "readySelector": ".app-loaded"
    },
    {
        "id": "dashboard_with_hover",
        "label": "Dashboard/Tooltip on Hover",
        "url": "http://localhost:3000/dashboard",
        "readySelector": ".dashboard-ready",
        "hoverSelector": ".chart-bar:first-child",
        "postInteractionWait": 500
    }
]
```

#### Scenario Options

| Option                | Type             | Description                                 |
|-----------------------|------------------|---------------------------------------------|
| `id`                  | string           | Unique identifier for the scenario          |
| `label`               | string           | Human-readable label (used in reports)      |
| `url`                 | string           | URL to navigate to                          |
| `browsers`            | string[]         | Override global browsers for this scenario  |
| `viewports`           | Viewport[]       | Override global viewports for this scenario |
| `readySelector`       | string           | CSS selector to wait for before capture     |
| `reloadAfterReady`    | boolean          | Reload page after ready selector appears    |
| `delay`               | number \| object | Wait time after ready (see below)           |
| `keyPressSelector`    | object           | Element to focus and key to press           |
| `hoverSelector`       | string           | Single element to hover over                |
| `hoverSelectors`      | array            | Multiple elements to hover in sequence      |
| `clickSelector`       | string           | Single element to click                     |
| `clickSelectors`      | array            | Multiple elements to click in sequence      |
| `postInteractionWait` | string \| number | Wait after interactions (selector or ms)    |
| `scrollToSelector`    | string           | Element to scroll into view                 |
| `misMatchThreshold`   | number           | Allowed mismatch percentage (0-100)         |

## Scenario Examples

### Basic Screenshot

Simple page capture waiting for a ready indicator:

```json
{
    "id": "simple_page",
    "label": "Simple Page",
    "url": "http://localhost:8080/page",
    "readySelector": ".page-loaded"
}
```

### With Custom Viewport

Override global viewports for specific scenarios:

```json
{
    "id": "wide_chart",
    "label": "Charts/Wide Chart View",
    "url": "http://localhost:8080/charts",
    "readySelector": ".chart-rendered",
    "viewports": [
        {
            "label": "wide",
            "width": 1920,
            "height": 1080
        }
    ]
}
```

### With Browser Override

Run a scenario only on specific browsers:

```json
{
    "id": "firefox_only",
    "label": "Firefox-specific Test",
    "url": "http://localhost:8080/test",
    "browsers": ["firefox"],
    "readySelector": ".ready"
}
```

### With Hover Interaction

Capture tooltip or hover state:

```json
{
    "id": "tooltip_test",
    "label": "Tooltip/Chart Hover",
    "url": "http://localhost:8080/chart",
    "readySelector": ".chart-ready",
    "hoverSelector": ".data-point:nth-child(3)",
    "postInteractionWait": 300
}
```

### With Multiple Hovers

Hover over multiple elements in sequence:

```json
{
    "id": "multi_hover",
    "label": "Multiple Hovers",
    "url": "http://localhost:8080/dashboard",
    "readySelector": ".loaded",
    "hoverSelectors": [
        ".menu-item:first-child",
        200,
        ".submenu-item"
    ],
    "postInteractionWait": ".tooltip-visible"
}
```

The `hoverSelectors` array supports:
- **Strings**: CSS selectors to hover
- **Numbers**: Milliseconds to wait before the next hover

### With Click Interaction

Capture state after clicking:

```json
{
    "id": "dropdown_open",
    "label": "Dropdown/Open State",
    "url": "http://localhost:8080/form",
    "readySelector": ".form-ready",
    "clickSelector": ".dropdown-toggle",
    "postInteractionWait": ".dropdown-menu"
}
```

### With Multiple Clicks

Click multiple elements in sequence:

```json
{
    "id": "wizard_step3",
    "label": "Wizard/Step 3",
    "url": "http://localhost:8080/wizard",
    "readySelector": ".wizard-loaded",
    "clickSelectors": [
        ".next-button",
        500,
        ".next-button",
        500,
        ".next-button"
    ],
    "postInteractionWait": ".step-3-content"
}
```

### With Key Press

Simulate keyboard input:

```json
{
    "id": "search_results",
    "label": "Search/Results View",
    "url": "http://localhost:8080/search",
    "readySelector": ".search-ready",
    "keyPressSelector": {
        "selector": ".search-input",
        "keyPress": "test query"
    },
    "postInteractionWait": ".results-loaded"
}
```

Supported special keys: `Enter`, `Tab`, `Backspace`, `Delete`, `Escape`, `ArrowUp`, `ArrowDown`, `ArrowLeft`, `ArrowRight`, `Home`, `End`, `PageUp`, `PageDown`, `F1`-`F12`, `Control`, `Alt`, `Meta`, `Shift`

Key combinations are also supported: `Control+a`, `Shift+Tab`

### With Scroll

Scroll to a specific element before capture:

```json
{
    "id": "footer_section",
    "label": "Page/Footer",
    "url": "http://localhost:8080/long-page",
    "readySelector": ".page-loaded",
    "scrollToSelector": ".footer-section",
    "postInteractionWait": 200
}
```

### With Delay

Add delays for animations or async content:

```json
{
    "id": "animated_chart",
    "label": "Charts/Animated",
    "url": "http://localhost:8080/animated-chart",
    "readySelector": ".chart-container",
    "delay": {
        "postReady": 1000,
        "postOperation": 500
    }
}
```

Simple delay (applied after ready):

```json
{
    "id": "simple_delay",
    "label": "With Delay",
    "url": "http://localhost:8080/page",
    "readySelector": ".ready",
    "delay": 2000
}
```

### With Mismatch Threshold

Allow small differences (useful for anti-aliasing variations):

```json
{
    "id": "chart_with_threshold",
    "label": "Charts/Allowable Variance",
    "url": "http://localhost:8080/chart",
    "readySelector": ".chart-ready",
    "misMatchThreshold": 0.25
}
```

### With Post-Interaction Wait

Wait for a selector after interactions:

```json
{
    "id": "async_content",
    "label": "Async/Content Load",
    "url": "http://localhost:8080/async",
    "readySelector": ".page-ready",
    "clickSelector": ".load-button",
    "postInteractionWait": ".content-loaded"
}
```

Or wait for a fixed duration:

```json
{
    "id": "animation_complete",
    "label": "Animation/Complete",
    "url": "http://localhost:8080/animation",
    "readySelector": ".ready",
    "clickSelector": ".animate-button",
    "postInteractionWait": 1500
}
```

### Complete Complex Example

```json
{
    "id": "complex_interaction",
    "label": "Dashboard/Full Interaction Flow",
    "url": "http://localhost:8080/dashboard",
    "browsers": ["chromium"],
    "viewports": [
        {
            "label": "hd",
            "width": 1920,
            "height": 1080
        }
    ],
    "readySelector": ".dashboard-ready",
    "delay": {
        "postReady": 500,
        "postOperation": 200
    },
    "clickSelectors": [
        ".filter-dropdown",
        300,
        ".filter-option:nth-child(2)"
    ],
    "hoverSelector": ".chart-bar:first-child",
    "postInteractionWait": ".tooltip-visible",
    "misMatchThreshold": 0.1
}
```

## Output Structure

After running tests, NeoBackstop generates the following output:

```
output/
├── reference/           # Reference screenshots (from approve mode)
│   └── *.png
├── test/                # Test screenshots (from test mode)
│   ├── *.png           # Current screenshots
│   └── diff_*.png      # Diff images for failures
├── html-report/         # Visual HTML report
│   ├── index.html
│   └── config.js
└── ci-report/           # Machine-readable results
    └── results.json
```

### Exit Codes

- `0`: All tests passed
- `1`: One or more tests failed or encountered errors

## Scenario Execution Order

Operations are executed in this order:

1. Navigate to URL (wait for `networkidle`)
2. Wait for `readySelector`
3. Reload page (if `reloadAfterReady` is true)
4. Apply `delay.postReady`
5. Execute `keyPressSelector`
6. Execute `hoverSelector`
7. Execute `hoverSelectors` (in order)
8. Execute `clickSelector`
9. Execute `clickSelectors` (in order)
10. Execute `scrollToSelector`
11. Apply `delay.postOperation`
12. Capture screenshot

## License

See [LICENSE](LICENSE) file for details.
