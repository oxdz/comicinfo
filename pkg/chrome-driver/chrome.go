// Package chromedriver is used to drive chrome browser
package chromedriver

import (
	"context"
	"net/http"
	"os"
	"path"
	"sync"

	"github.com/chromedp/cdproto/network"
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

	chromedp.Flag("disable-gpu", true),
	chromedp.Flag("enable-automation", false),
	chromedp.Flag("excludeSwitches", "enable-automation"),
}

func Dir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "."
	}
	return path.Join(homeDir, "ccinf", ".chromedp")
}

func NewExecAllocatorCtx(dir string, headless bool) (context.Context, context.CancelFunc) {
	opts := append(DefaultExecAllocatorOptions[:],
		chromedp.DisableGPU,
		chromedp.Flag("headless", headless))

	if dir != "" {
		opts = append(opts, chromedp.UserDataDir(dir))
	}

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts[:]...)

	return allocCtx, cancel
}

func NewChromeContext(headless bool) (context.Context, context.CancelFunc) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", headless), // 禁用无头模式
		chromedp.Flag("disable-gpu", true),
		// chromedp.Flag("no-sandbox", true),
		chromedp.Flag("enable-automation", false),
		// chromedp.Flag("disable-blink-features", "AutomationControlled"),
		chromedp.Flag("excludeSwitches", "enable-automation"),
	)

	allocCtx, cancelAlloc := chromedp.NewExecAllocator(
		context.Background(), opts...)

	// create chrome instance
	ctx, cancel := chromedp.NewContext(
		allocCtx,
		// chromedp.WithDebugf(log.Printf),
	)

	return ctx, func() { cancel(); cancelAlloc() }
}

type Cookies struct {
	Domain map[string]int
	KVs    []map[[2]string]*http.Cookie
	sync.Mutex
}

func (c *Cookies) Add(v network.Cookie) {
	c.Lock()
	defer c.Unlock()
	if c.Domain == nil {
		c.Domain = make(map[string]int)
	}

	if _, ok := c.Domain[v.Domain]; !ok {
		c.Domain[v.Domain] = len(c.KVs)
		c.KVs = append(c.KVs, map[[2]string]*http.Cookie{{
			v.Name, v.Path,
		}: &http.Cookie{
			Name:   v.Name,
			Value:  v.Value,
			Path:   v.Path,
			Domain: v.Domain,
			Secure: v.Secure,
		}})
	} else {
		c.KVs[c.Domain[v.Domain]][[2]string{
			v.Name, v.Path}] = &http.Cookie{
			Name:   v.Name,
			Value:  v.Value,
			Path:   v.Path,
			Domain: v.Domain,
			Secure: v.Secure,
		}
	}
}

func (c *Cookies) ActionFunc() chromedp.ActionFunc {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		cookies, err := network.GetCookies().Do(ctx)
		if err != nil {
			return err
		}
		for _, cookie := range cookies {
			c.Add(*cookie)
		}
		return nil
	})
}

func (c *Cookies) GetAll(domain ...string) (r []*http.Cookie) {
	c.Lock()
	defer c.Unlock()
	for _, domain := range domain {
		if v, ok := c.Domain[domain]; ok {
			for _, v := range c.KVs[v] {
				r = append(r, v)
			}
		}
	}
	return r
}
