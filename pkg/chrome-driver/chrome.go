// Package chromedriver is used to drive chrome browser
package chromedriver

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/chromedp/chromedp"
)

var DefaultExecAllocatorOptions = [...]chromedp.ExecAllocatorOption{
	chromedp.NoFirstRun,
	chromedp.NoDefaultBrowserCheck,
	chromedp.Flag("disable-background-networking", true),
	chromedp.Flag("enable-features", "NetworkService,NetworkServiceInProcess"),
	chromedp.Flag("disable-background-timer-throttling", true),
	chromedp.Flag("disable-backgrounding-occluded-windows", true),
	chromedp.Flag("disable-breakpad", true),
	chromedp.Flag("disable-client-side-phishing-detection", true),
	chromedp.Flag("disable-default-apps", true),
	chromedp.Flag("disable-dev-shm-usage", true),
	chromedp.Flag("disable-extensions", true),
	chromedp.Flag("disable-features", "site-per-process,Translate,BlinkGenPropertyTrees"),
	chromedp.Flag("hide-crash-restore-bubble", true),
	chromedp.Flag("disable-hang-monitor", true),
	chromedp.Flag("disable-ipc-flooding-protection", true),
	chromedp.Flag("disable-popup-blocking", true),
	chromedp.Flag("disable-prompt-on-repost", true),
	chromedp.Flag("disable-renderer-backgrounding", true),
	chromedp.Flag("disable-sync", true),
	chromedp.Flag("force-color-profile", "srgb"),
	chromedp.Flag("metrics-recording-only", true),
	chromedp.Flag("safebrowsing-disable-auto-update", true),
	chromedp.Flag("password-store", "basic"),
	chromedp.Flag("use-mock-keychain", true),

	chromedp.Flag("enable-automation", false),
	// chromedp.Flag("disable-blink-features", "AutomationControlled"),
}

func NewExecAllocatorCtx(dir string) (context.Context, context.CancelFunc, error) {
	if dir == "" {
		dir = "./ccinf-web-tmp-cache"
	}
	if fsinf, err := os.Stat(dir); err != nil {
		// log.Print(err)
		_ = os.MkdirAll(dir, 0o755)
	} else {
		if !fsinf.IsDir() {
			return nil, nil, fmt.Errorf("%s is not a directory", dir)
		}
	}

	dir = filepath.Join(dir, "userdatadir")
	if err := os.Mkdir(dir, 0o755); err != nil {
		if !os.IsExist(err) {
			return nil, nil, fmt.Errorf("mkdir userdatadir err: %+v", err)
		}
	}

	opts := append(DefaultExecAllocatorOptions[:],
		chromedp.DisableGPU,
		chromedp.UserDataDir(dir),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts[:]...)

	// also set up a custom logger

	return allocCtx, cancel, nil
}
