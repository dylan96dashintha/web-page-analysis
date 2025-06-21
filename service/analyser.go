package service

import (
	"context"
	"github.com/PuerkitoBio/goquery"
	"github.com/web-page-analysis/container"
	"github.com/web-page-analysis/domain"
	"github.com/web-page-analysis/usecase"
	"net/http"
)

type Analyser interface {
	WebAnalyser(ctx context.Context, req domain.AnalyserRequest)
}

type analyser struct {
	container container.Container
}

func NewAnalyser(ctr container.Container) Analyser {
	return &analyser{
		container: ctr,
	}
}

func (a analyser) WebAnalyser(ctx context.Context, req domain.AnalyserRequest) {

	analyserObj := usecase.NewAnalyser(a.container)
	for _, url := range req.Url {
		resp, err := a.container.OBAdapter.Get(ctx, url)
		if err != nil {
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return
		}

		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			return
		}
		analyserObj.GetTitle(ctx, doc)

	}

}
