package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"testing"

	"github.com/chromedp/chromedp"
	chromedriver "github.com/oxdz/comicinfo/pkg/chrome-driver"
)

func TestCmd(t *testing.T) {
	dir := chromedriver.Dir()
	fmt.Printf("chrome user data saved to dir: `%s`\n", dir)
	ctx, alloccancel := chromedriver.NewExecAllocatorCtx(dir, false)
	defer alloccancel()

	ctx, cancel := chromedp.NewContext(ctx)
	defer cancel()

	if err := chromedp.Run(ctx); err != nil {
		t.Fatal(err)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTSTP)
	fmt.Printf("Press Ctrl+Z to exit\n")

	ck := &chromedriver.Cookies{}
	err := sh(ctx, ck, "https://pocket.shonenmagazine.com/title/01774/episode/349270", "/tmp/aaaazax")
	if err != nil {
		t.Fatal(err)
	}
}
