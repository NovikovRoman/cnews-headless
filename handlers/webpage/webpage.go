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
	body    []byte
	cookies []string
}

func New(u string) (w *Webpage) {
	w = &Webpage{
		url: u,
	}
	return
}

func (w *Webpage) String() string {
	return string(w.body)
}

func (w *Webpage) Bytes() []byte {
	return w.body
}

func (w *Webpage) Cookies() []string {
	return w.cookies
}

func (w *Webpage) File(ctx context.Context, cookies ...string) (err error) {
	w.cookies = cookies
	opts := []chromedp.ExecAllocatorOption{
		chromedp.NoFirstRun,
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

	done := make(chan bool)
	var requestID network.RequestID
	chromedp.ListenTarget(ctx, func(v interface{}) {
		switch ev := v.(type) {
		case *network.EventRequestWillBeSent:
			log.Infof("EventRequestWillBeSent: %v: %v", ev.RequestID, ev.Request.URL)
			if ev.Request.URL == w.url {
				requestID = ev.RequestID
			}
		case *network.EventLoadingFinished:
			log.Infof("EventLoadingFinished: %v", ev.RequestID)
			if ev.RequestID == requestID {
				close(done)
			}
		}
	})

	start := time.Now()
	err = chromedp.Run(
		ctx,
		network.Enable(),
		w.setCookies(),
		chromedp.Navigate(w.url),
		w.getCookies(),
	)
	log.Infof("%s %f sec", w.url, time.Since(start).Seconds())
	if err != nil {
		return
	}
	<-done

	err = chromedp.Run(ctx, chromedp.ActionFunc(func(ctx context.Context) (err error) {
		w.body, err = network.GetResponseBody(requestID).Do(ctx)
		return
	}))
	return
}

func (w *Webpage) Html(ctx context.Context, selector string, cookies ...string) (err error) {
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

	content := ""
	err = chromedp.Run(
		ctx,
		network.Enable(),
		w.setCookies(),
		chromedp.Navigate(w.url),
		chromedp.WaitReady(selector),
		w.getCookies(),
		chromedp.OuterHTML("html", &content, chromedp.ByQuery),
	)
	log.Infof("%s %f sec", w.url, time.Since(start).Seconds())
	if err != nil {
		return
	}

	w.body = []byte(html.UnescapeString(content))
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
