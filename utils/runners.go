package utils
import(
    "context"
    "math"
    "sync"
	"log"
	"errors"
    "github.com/chromedp/cdproto/emulation"
    "github.com/chromedp/cdproto/page"
    "github.com/chromedp/chromedp"
)

func Main_runner(url string, wg *sync.WaitGroup, ch chan<- *SSResponse) {
    // create context
    ctx, cancel := chromedp.NewContext(context.Background())
    defer cancel()
    defer wg.Done()
    // capture screenshot of an element
    var buf []byte
	var err error
	valid,isok :=HTTPify(url)
	if isok == true{
    // capture entire browser viewport, returning png with quality=90
		if err = chromedp.Run(ctx, fullScreenshot(valid, 90, &buf)); err != nil {
			log.Println(err)
		}
	}else{
		log.Println("invalid URL "+url)
		err = errors.New("Invalid URL")
	}
	ch <- &SSResponse{url,err,buf}
}

// fullScreenshot takes a screenshot of the entire browser viewport.

func fullScreenshot(urlstr string, quality int64, res *[]byte) chromedp.Tasks {
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

