package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
	"github.com/web-page-analysis/bootstrap"
	"github.com/web-page-analysis/container"
	"github.com/web-page-analysis/domain"
	"github.com/web-page-analysis/usecase"
	"io"
	"net/http"
	"strings"
	"sync"
)

const (
	prefix = "service.analyser "
)

type Analyser interface {
	WebAnalyser(ctx context.Context, req domain.AnalyserRequest) (res domain.AnalysisResult, err error)
}

type analyser struct {
	container container.Container
	config    bootstrap.Config
}

func NewAnalyser(ctr container.Container, config bootstrap.Config) Analyser {
	return &analyser{
		container: ctr,
		config:    config,
	}
}

func (a analyser) WebAnalyser(ctx context.Context, req domain.AnalyserRequest) (res domain.AnalysisResult, err error) {
	log.WithContext(ctx).Info(prefix, "start to analyse the url")
	// validate the url
	validatorObj := usecase.NewValidation()
	isValid := validatorObj.IsValidUrl(ctx, req.Url)
	if !isValid {
		log.WithContext(ctx).Error(prefix, "Invalid url")
		return res, errors.New("invalid url")
	}

	// start analysing the webpage
	analyserObj := usecase.NewAnalyser(a.container, a.config)

	// call the webpage to get the html
	resp, err := a.container.OBAdapter.Get(ctx, req.Url)
	if err != nil {
		log.WithContext(ctx).Error(prefix, "Error in calling outbound call, err: ", err)
		return res, err
	}
	if resp != nil && resp.StatusCode != http.StatusOK {
		log.WithContext(ctx).Error(prefix, "Error in calling outbound call, status: ", resp.StatusCode)
		return res, errors.New(fmt.Sprintf("Error in reaching server,  status: %s", resp.Status))
	}

	// close the resp body in need to close the file descriptor in resource level
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.WithContext(ctx).Error(prefix, "Error in reading response body, err: ", err)
		return res, err
	}
	bodyString := string(bodyBytes)
	// need to read the resp.Body twice, to overcome this,used this technique
	resp.Body = io.NopCloser(strings.NewReader(bodyString))

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.WithContext(ctx).Error(prefix, "Data cannot be parsed to HTML, err: ", err)
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

	// get the title of the html
	wg.Add(1)
	go func() {
		defer wg.Done()
		title = analyserObj.GetTitle(ctx, doc)
		return
	}()

	// get the html version of the html
	wg.Add(1)
	go func() {
		defer wg.Done()
		htmlVersion = analyserObj.CheckHtmlVersion(ctx, bodyString)
		return
	}()

	// check any logins are there in the html
	wg.Add(1)
	go func() {
		defer wg.Done()
		login = analyserObj.CheckAnyLogin(ctx, doc)
		return

	}()

	// check any links are there in the html
	wg.Add(1)
	go func() {
		defer wg.Done()
		link = analyserObj.CountLinks(ctx, doc, req.Url)
		return

	}()

	// count the heading types in the html
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
