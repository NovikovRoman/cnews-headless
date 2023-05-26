package hlshell

import (
	"context"
	"fmt"
	"html"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	log "github.com/sirupsen/logrus"
)

type HeadlessShell struct {
	ctx     context.Context
	cookies []string
}

func New(ctx context.Context) (w *HeadlessShell) {
	w = &HeadlessShell{
		ctx: ctx,
	}
	return
}

func (h *HeadlessShell) Cookies() []string {
	return h.cookies
}

func (h *HeadlessShell) File(url string) (body []byte, err error) {
	ctx, cancel := context.WithTimeout(h.ctx, time.Second*45)
	defer cancel()

	done := make(chan bool)
	var requestID network.RequestID
	chromedp.ListenTarget(ctx, func(v interface{}) {
		switch ev := v.(type) {
		case *network.EventRequestWillBeSent:
			log.Infof("EventRequestWillBeSent: %v: %v", ev.RequestID, ev.Request.URL)
			if ev.Request.URL == url {
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
		chromedp.Navigate(url),
	)
	log.Infof("%s %f sec", url, time.Since(start).Seconds())
	if err != nil {
		return
	}
	<-done

	err = chromedp.Run(ctx, chromedp.ActionFunc(func(ctx context.Context) (err error) {
		body, err = network.GetResponseBody(requestID).Do(ctx)
		return
	}))
	return
}

func (h *HeadlessShell) Html(url, selector string) (content string, err error) {
	ctx, cancel := context.WithTimeout(h.ctx, time.Second*45)
	defer cancel()

	if selector == "" {
		selector = "body"
	}
	start := time.Now()

	err = chromedp.Run(
		ctx,
		network.Enable(),
		chromedp.Navigate(url),
		chromedp.WaitReady(selector),
		h.getCookies(),
		chromedp.OuterHTML("html", &content, chromedp.ByQuery),
	)
	log.Infof("%s %f sec", url, time.Since(start).Seconds())
	if err != nil {
		return
	}

	content = html.UnescapeString(content)
	return
}

func (h *HeadlessShell) getCookies() chromedp.Action {
	return chromedp.ActionFunc(func(ctx context.Context) (err error) {
		var cookies []*network.Cookie
		if cookies, err = network.GetCookies().Do(ctx); err != nil {
			return
		}

		h.cookies = make([]string, len(cookies))
		for i, cookie := range cookies {
			h.cookies[i] = fmt.Sprintf("%s=%s", cookie.Name, cookie.Value)
		}
		return
	})
}
