package commands

import (
	"context"
	"io/ioutil"
	"math"
	"strings"

	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/florinutz/filme/pkg/filme"
	"github.com/spf13/cobra"
)

// BuildScreenshotCmd is only available for debug because it's only useful for
// checking whether sites are available to headless browsing (crawlers) or not.
//
// https://github.com/chromedp/examples/blob/master/screenshot/main.go
func BuildScreenshotCmd(f *filme.Filme) *cobra.Command {
	var opts struct {
		jpegPath string
	}

	cmd := &cobra.Command{
		Use:   "screenshot <url>",
		Short: "takes a url's screenshot",

		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := chromedp.NewContext(context.Background())
			defer cancel()

			var buf []byte

			if err := chromedp.Run(ctx, screenshot(strings.Join(args, " "), 90, &buf)); err != nil {
				f.Log.Fatal(err)
			}
			if err := ioutil.WriteFile(opts.jpegPath, buf, 0644); err != nil {
				f.Log.Fatal(err)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&opts.jpegPath, "jpeg-path", "p", "screenshot.jpg",
		"output path")

	return cmd
}

func screenshot(urlstr string, quality int64, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.ActionFunc(func(ctx context.Context) error {
			// get layout metrics
			_, _, contentSize, err := page.GetLayoutMetrics().Do(ctx)
			if err != nil {
				return err
			}

			width, height := int64(math.Ceil(contentSize.Width)), int64(math.Ceil(contentSize.Height))

			// force viewport emulation
			err = emulation.SetDeviceMetricsOverride(width, height, 1, false).
				WithScreenOrientation(&emulation.ScreenOrientation{
					Type:  emulation.OrientationTypePortraitPrimary,
					Angle: 0,
				}).
				Do(ctx)
			if err != nil {
				return err
			}

			// capture screenshot
			*res, err = page.CaptureScreenshot().
				WithQuality(quality).
				WithClip(&page.Viewport{
					X:      contentSize.X,
					Y:      contentSize.Y,
					Width:  contentSize.Width,
					Height: contentSize.Height,
					Scale:  1,
				}).Do(ctx)
			if err != nil {
				return err
			}
			return nil
		}),
	}
}
