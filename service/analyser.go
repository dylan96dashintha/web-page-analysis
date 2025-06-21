package service

import (
	"context"
	"errors"
	"github.com/PuerkitoBio/goquery"
	"github.com/web-page-analysis/container"
	"github.com/web-page-analysis/domain"
	"github.com/web-page-analysis/usecase"
	"io"
	"net/http"
	"strings"
	"sync"
)

type Analyser interface {
	WebAnalyser(ctx context.Context, req domain.AnalyserRequest) (res domain.AnalysisResult, err error)
}

type analyser struct {
	container container.Container
}

func NewAnalyser(ctr container.Container) Analyser {
	return &analyser{
		container: ctr,
	}
}

func (a analyser) WebAnalyser(ctx context.Context, req domain.AnalyserRequest) (res domain.AnalysisResult, err error) {

	analyserObj := usecase.NewAnalyser(a.container)
	resp, err := a.container.OBAdapter.Get(ctx, req.Url)
	if err != nil {
		return res, err
	}
	if resp != nil && resp.StatusCode != http.StatusOK {
		return res, errors.New("Unable in reaching to the server")
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return res, err
	}
	bodyString := string(bodyBytes)
	resp.Body = io.NopCloser(strings.NewReader(bodyString))

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return res, err
	}

	wg := new(sync.WaitGroup)
	var (
		title       string
		htmlVersion string
		login       bool
		link        domain.Link
		heading     map[string]int
	)
	wg.Add(1)
	go func() {
		defer wg.Done()
		title = analyserObj.GetTitle(ctx, doc)
		return
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		htmlVersion = analyserObj.CheckHtmlVersion(ctx, bodyString)
		return
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		login = analyserObj.CheckAnyLogin(ctx, doc)
		return

	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		link = analyserObj.CountLinks(ctx, doc, req.Url)
		return

	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		heading = analyserObj.CountHeading(ctx, doc)
		return
	}()

	wg.Wait()
	result := domain.AnalysisResult{
		HTMLVersion:  htmlVersion,
		Title:        title,
		Headings:     heading,
		Link:         link,
		HasLoginForm: login,
	}

	return result, nil

}
