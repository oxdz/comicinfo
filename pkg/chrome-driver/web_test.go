package chromedriver

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/chromedp/chromedp"
)

func TestNewTaskCtx(t *testing.T) {
	ctx, alloccancel, err := NewExecAllocatorCtx("")
	if err != nil {
		t.Fatal(err)
	}
	defer alloccancel()

	ctx, cancel := chromedp.NewContext(ctx)

	if err = chromedp.Run(ctx); err != nil {
		t.Fatal(err)
	}

	ctxWithTimeout, cancelWt := context.WithTimeout(ctx, time.Second*2)
	defer cancelWt()

	if err = chromedp.Run(ctxWithTimeout, chromedp.Navigate("https://www.google.com")); err != nil {
		if errors.Is(context.DeadlineExceeded, err) {
			t.Log("timeout")
		} else {
			t.Fatal(err)
		}
	}

	cancel()
}
