package webpage

import (
	"context"
	"html"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	log "github.com/sirupsen/logrus"
)

type Webpage struct {
	url     string
	html    string
	cookies []string
}

func New(u string) (w *Webpage) {
	w = &Webpage{
		url: u,
	}
	return
}

func (w *Webpage) String() string {
	return w.html
}

func (w *Webpage) Cookies() []string {
	return w.cookies
}

func (w *Webpage) Get(ctx context.Context, selector string, cookies ...string) (err error) {
	w.cookies = cookies
	opts := []chromedp.ExecAllocatorOption{
		// chromedp.ExecPath("/headless-shell"),
		//chromedp.UserAgent("Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.5672.127 Safari/537.36"),
		//chromedp.WindowSize(1280, 720),
		chromedp.NoFirstRun,
		//chromedp.NoDefaultBrowserCheck,
		chromedp.Headless,
		chromedp.IgnoreCertErrors,
		chromedp.DisableGPU,
	}

	ctx, cancel := chromedp.NewExecAllocator(ctx, opts...)
	defer cancel()

	ctx, cancel = chromedp.NewContext(ctx, chromedp.WithLogf(log.Printf))
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 45*time.Second)
	defer cancel()

	if selector == "" {
		selector = "body"
	}
	start := time.Now()
	err = chromedp.Run(
		ctx,
		network.Enable(),
		w.setCookies(),
		chromedp.Navigate(w.url),
		chromedp.WaitReady(selector),
		w.getCookies(),
		chromedp.OuterHTML("html", &w.html, chromedp.ByQuery),
		// chromedp.WaitNotVisible(`#trk_jschal_nojs`, chromedp.ByQuery),
		// cookie(),
		//chromedp.FullScreenshot(&b, 80),
		//removeCookie(),
	)
	log.Infof("%s %f sec", w.url, time.Since(start).Seconds())
	if err != nil {
		return
	}

	w.html = html.UnescapeString(w.html)
	return
}

func (w *Webpage) setCookies() chromedp.Action {
	return chromedp.ActionFunc(func(ctx context.Context) (err error) {
		if len(w.cookies) == 0 {
			log.Info("set empty cookies")
			return
		}
		var cookies []*network.CookieParam
		for _, cookie := range w.cookies {
			c := &network.CookieParam{}
			if err = c.UnmarshalJSON([]byte(cookie)); err != nil {
				log.Errorf("%s %v", cookie, err)
				continue
			}
			cookies = append(cookies, c)
		}
		network.SetCookies(cookies)
		return
	})
}

func (w *Webpage) getCookies() chromedp.Action {
	return chromedp.ActionFunc(func(ctx context.Context) (err error) {
		var cookies []*network.Cookie
		if cookies, err = network.GetCookies().Do(ctx); err != nil {
			return
		}

		w.cookies = make([]string, len(cookies))
		for i, cookie := range cookies {
			expr := cdp.TimeSinceEpoch(time.Unix(int64(cookie.Expires), 0))
			param := &network.CookieParam{
				Name:         cookie.Name,
				Value:        cookie.Value,
				URL:          w.url,
				Domain:       cookie.Domain,
				Path:         cookie.Path,
				Secure:       cookie.Secure,
				HTTPOnly:     cookie.HTTPOnly,
				SameSite:     cookie.SameSite,
				Expires:      &expr,
				Priority:     cookie.Priority,
				SameParty:    cookie.SameParty,
				SourceScheme: cookie.SourceScheme,
				SourcePort:   cookie.SourcePort,
				PartitionKey: cookie.PartitionKey,
			}
			var b []byte
			if b, err = param.MarshalJSON(); err != nil {
				continue
			}
			w.cookies[i] = string(b)
		}
		return
	})
}
