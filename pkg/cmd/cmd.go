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
	"os/signal"
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
	fmt.Printf("Chrome user data saved to dir: `%s`\n", dir)
	ctx, alloccancel := chromedriver.NewExecAllocatorCtx(dir, false)
	defer alloccancel()

	ctx, cancel := chromedp.NewContext(ctx)
	defer cancel()

	if err := chromedp.Run(ctx); err != nil {
		return err
	}

	fmt.Printf("Press Ctrl+C to exit\n\n")
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	ck := &chromedriver.Cookies{}
	lnCh, errCh := Scanner(os.Stdin)

	homedir, err := os.UserHomeDir()
	if err != nil {
		homedir = "."
	}
	baseDir := path.Join(homedir, "ccinf")
	imageBaseDir := path.Join(baseDir, "images")

	for {
		fmt.Printf("url: ")
		select {
		case err := <-errCh:
			return err
		case <-sig:
			fmt.Printf("\nExit!\n")
			return nil
		case line := <-lnCh:
			if line == "" {
				continue
			}
			switch {
			case shonenmagazineRe.MatchString(line):
				if err := sh(ctx, ck, line, imageBaseDir); err != nil {
					fmt.Println(err.Error())
				}
				fmt.Printf("Done!\n")
			}
		}
	}
}

func Scanner(r io.Reader) (<-chan string, <-chan error) {
	ch := make(chan string)
	errCh := make(chan error, 1)
	scanner := bufio.NewScanner(r)
	go func() {
		defer close(errCh)
		defer close(ch)

		for scanner.Scan() {
			line := scanner.Text()
			line = strings.TrimSpace(line)
			ch <- line
		}

		err := scanner.Err()
		if err != nil {
			err = fmt.Errorf("scanner error: %w", err)
		}
		errCh <- err
	}()
	return ch, errCh
}

func sh(ctx context.Context, ck *chromedriver.Cookies, url, basedir string) error {
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

	dst := path.Join(basedir, m[1], m[2])
	os.MkdirAll(dst, 0755)
	fmt.Printf("save images to: `%s`\n", dst)

	for i, url := range inf.PageList {
		bar(i, len(inf.PageList))
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
	bar(len(inf.PageList), len(inf.PageList))
	return nil
}

func bar(currentPage, totalPages int) {
	barLength := 50

	percent := float64(currentPage) / float64(totalPages) * 100

	filledLength := int(percent / 100 * float64(barLength))

	bar := "["
	for i := 0; i < barLength; i++ {
		if i < filledLength {
			bar += "="
		} else if i == filledLength {
			bar += ">"
		} else {
			bar += " "
		}
	}
	bar += "]"

	fmt.Printf("\r%s pages: %d/%d (%.1f%%)", bar, currentPage, totalPages, percent)
	if currentPage == totalPages {
		fmt.Printf("\n")
	}
}
