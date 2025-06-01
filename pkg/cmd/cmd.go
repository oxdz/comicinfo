package cmd

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"image"
	"io"
	"net/http"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/oxdz/comicinfo/internal/reader/shonenmagazine"
	chromedriver "github.com/oxdz/comicinfo/pkg/chrome-driver"
	"github.com/oxdz/comicinfo/pkg/decode"
	"github.com/spf13/cobra"
)

var (
	shonenmagazineRe = regexp.MustCompile(`^https://pocket.shonenmagazine.com/title/(\d+)/episode/(\d+)`)
)

func Start(cmd *cobra.Command, args []string) error {
	dir := chromedriver.Dir()
	fmt.Printf("chrome user data saved to dir: `%s`\n", dir)
	ctx, alloccancel := chromedriver.NewExecAllocatorCtx(dir, false)
	defer alloccancel()

	ctx, cancel := chromedp.NewContext(ctx)
	defer cancel()

	if err := chromedp.Run(ctx); err != nil {
		return err
	}

	fmt.Printf("Press Ctrl+Z to exit\n")

	ck := &chromedriver.Cookies{}
	for {
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Printf("Please input url:\n")
		for scanner.Scan() {
			line := scanner.Text()
			line = strings.TrimSpace(line)
			switch {
			case shonenmagazineRe.MatchString(line):
				if err := sh(ctx, ck, line); err != nil {
					fmt.Println(err.Error())
				}
			}
		}

		// 检查扫描过程中是否有错误
		if err := scanner.Err(); err != nil {
			fmt.Printf("scanner error: %v\n", err)
		}
	}
}

func sh(ctx context.Context, ck *chromedriver.Cookies, url string) error {
	m := shonenmagazineRe.FindStringSubmatch(url)
	var err error
	for range 3 {
		if err = chromedp.Run(ctx,
			chromedp.Navigate(url),
			ck.ActionFunc()); err != nil {
			fmt.Println(err)
			time.Sleep(time.Second)
		} else {
			break
		}
	}

	if err != nil {
		return fmt.Errorf("failed to open url(%s): %w", url, err)
	}

	var inf *decode.EpisodeInfo
	for range 3 {
		id, _ := strconv.Atoi(m[2])
		inf, err = shonenmagazine.EpisodeInfo(ctx, id,
			ck.GetAll(".shonenmagazine.com", ".pocket.shonenmagazine.com"))
		if err != nil {
			fmt.Println(err.Error())
			time.Sleep(time.Second)
		} else {
			break
		}
	}

	if err != nil {
		return fmt.Errorf("failed to get episode(%s) info: %w", m[2], err)
	}

	if inf == nil {
		return fmt.Errorf("failed to get episode(%s) info", m[2])
	}

	homedir, err := os.UserHomeDir()
	if err != nil {
		homedir = "."
	}

	dst := path.Join(homedir, "ccinf/images", m[1], m[2])
	os.MkdirAll(dst, 0755)
	fmt.Printf("save images to: %s\n", dst)

	for i, url := range inf.PageList {
		r, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return err
		}
		var resp *http.Response
		for range 3 {
			resp, err = http.DefaultClient.Do(r)
			if err != nil {
				fmt.Println(err.Error())
			} else {
				break
			}
		}
		if err != nil {
			return err
		}

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("failed to get response: status: %s", resp.Status)
		}

		defer resp.Body.Close()
		buf, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		tp := resp.Header.Get("Content-Type")
		switch tp {
		case "image/png":
		case "image/jpeg":
		default:
			return fmt.Errorf("unexpected content-type: %s", tp)
		}

		im, _, err := image.Decode(bytes.NewReader(buf))
		if err != nil {
			return err
		}

		if err := shonenmagazine.DrawImage(im,
			path.Join(dst, fmt.Sprintf("%d.png", i)), inf.ScrambleSeed); err != nil {
			return err
		}
	}

	return nil
}
