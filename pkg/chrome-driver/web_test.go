package chromedriver

import (
	"io"
	"net/http"
	"os"
	"path"
	"strconv"
	"testing"

	"github.com/chromedp/chromedp"
	"github.com/oxdz/comicinfo/internal/reader/shonenmagazine"
)

func TestNewTaskCtx(t *testing.T) {
	ctx, alloccancel := NewExecAllocatorCtx("/tmp/aaaazax", false)
	defer alloccancel()

	ctx, cancel := chromedp.NewContext(ctx)

	if err := chromedp.Run(ctx); err != nil {
		t.Fatal(err)
	}

	ck := &Cookies{}

	if err := chromedp.Run(ctx,
		chromedp.Navigate("https://pocket.shonenmagazine.com/title/01774/episode/349270"),
		ck.ActionFunc()); err != nil {
		t.Error(err)
	}

	// https://pocket.shonenmagazine.com/title/01774/episode/349272
	inf, err := shonenmagazine.EpisodeInfo(ctx, 349270,
		ck.GetAll(".shonenmagazine.com", ".pocket.shonenmagazine.com"))

	t.Log(inf, err)

	dst := "/tmp/images/349270"
	os.MkdirAll(dst, 0755)
	for i, url := range inf.PageList {
		r, err := http.NewRequest("GET", url, nil)
		if err != nil {
			t.Error(err)
		}
		resp, err := http.DefaultClient.Do(r)
		if err != nil {
			t.Error(err)
		}
		defer resp.Body.Close()
		buf, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Error(err)
		}
		tp := resp.Header.Get("Content-Type")
		switch tp {
		case "image/png":
			d := path.Join(dst, strconv.Itoa(i)+".png")
			os.WriteFile(d, buf, 0644)
		case "image/jpeg":
			d := path.Join(dst, strconv.Itoa(i)+".jpeg")
			os.WriteFile(d, buf, 0644)
		default:
			t.Error("unexpected content-type", t)
		}

	}

	// if err := chromedp.Run(ctx,
	// 	chromedp.Navigate("https://pocket.shonenmagazine.com/title/01774/episode/349272"),
	// 	ck.ActionFunc()); err != nil {
	// 	t.Error(err)
	// }

	// inf2, err := shonenmagazine.EpisodeInfo(ctx, 349272,
	// 	ck.GetAll(".shonenmagazine.com", ".pocket.shonenmagazine.com"))

	// t.Log(inf2, err)

	cancel()
	t.Error("")
}
