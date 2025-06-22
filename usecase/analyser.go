package usecase

import (
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
	"github.com/web-page-analysis/bootstrap"
	"github.com/web-page-analysis/container"
	"github.com/web-page-analysis/domain"
	"strings"
	"sync"
)

const (
	title          = "title"
	analyserPrefix = "usecase.analyser "
)

type Analyser interface {
	CheckHtmlVersion(ctx context.Context, rawHtml string) string
	GetTitle(ctx context.Context, doc *goquery.Document) string
	CountHeading(ctx context.Context, doc *goquery.Document) map[string]int
	CountLinks(ctx context.Context, doc *goquery.Document, url string) domain.Link
	CheckAnyLogin(ctx context.Context, doc *goquery.Document) bool
}

type analyser struct {
	ctr    container.Container
	config bootstrap.Config
}

// CheckHtmlVersion check the html version
// if the condition passed, then it returns
// otherwise it returns as unknown
func (a analyser) CheckHtmlVersion(ctx context.Context, rawHTML string) string {
	log.WithContext(ctx).Info(analyserPrefix, "start to checking HTML version")
	rawHTML = strings.ToLower(rawHTML)
	switch {
	case strings.Contains(rawHTML, `<!doctype html>`):
		return "HTML5"
	case strings.Contains(rawHTML, `<!doctype html public "-//w3c//dtd html 4.01"`):
		return "HTML 4.01"
	case strings.Contains(rawHTML, `<!doctype html public "-//w3c//dtd xhtml 1.0"`):
		return "XHTML 1.0"
	case strings.Contains(rawHTML, `<!doctype html public "-//w3c//dtd html 3.2"`):
		return "HTML 3.2"
	default:
		return "Unknown or missing doctype"
	}
}

func (a analyser) CheckAnyLogin(ctx context.Context, doc *goquery.Document) bool {
	log.WithContext(ctx).Info(analyserPrefix, "start to checking login")
	var isExist bool
	doc.Find("form").Each(func(i int, s *goquery.Selection) {
		s.Find("input").Each(func(j int, input *goquery.Selection) {
			inputType, _ := input.Attr("type")
			inputName, _ := input.Attr("name")

			if strings.ToLower(inputType) == "password" ||
				strings.ToLower(inputType) == "email" ||
				strings.ToLower(inputName) == "username" {
				isExist = true
			}
		})

		if !isExist {
			text := strings.ToLower(s.Text())
			if strings.Contains(text, "log in") || strings.Contains(text, "login") {
				isExist = true
			}
		}
	})

	return isExist
}

func (a analyser) GetTitle(ctx context.Context, doc *goquery.Document) string {
	log.WithContext(ctx).Info(analyserPrefix, "start to fetching the title")
	return doc.Find(title).Text()
}

func (a analyser) CountHeading(ctx context.Context, doc *goquery.Document) map[string]int {
	log.WithContext(ctx).Info(analyserPrefix, "start to counting the heading")
	headingsMap := map[string]int{}
	for i := 1; i <= 6; i++ {
		tag := fmt.Sprintf("h%d", i)
		headingsMap[tag] = doc.Find(tag).Length()
	}

	return headingsMap
}

func (a analyser) CountLinks(ctx context.Context, doc *goquery.Document, baseURL string) domain.Link {
	log.WithContext(ctx).Info(analyserPrefix, "start to counting the links")
	var (
		link          domain.Link
		linkMu        sync.Mutex
		linkJobs      = make(chan string)
		wg            sync.WaitGroup
		distinctLinks = make(map[string]interface{})
	)

	link.InaccessibleLink = make([]string, 0)

	baseURL = normalizeURL(baseURL)

	// initiate the worker pool with the config value
	// only to check the accessibility of the links
	for i := 0; i < int(a.config.AppConfig.WorkerCount); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for href := range linkJobs {
				var inaccessible bool
				fullURL := resolveURL(baseURL, href)
				resp, err := a.ctr.OBAdapter.Get(ctx, fullURL)
				if resp != nil && resp.Body != nil {
					defer resp.Body.Close()
				}

				if err != nil || resp != nil && (resp.StatusCode > 300 || resp.StatusCode < 200) {
					log.WithContext(ctx).Error(analyserPrefix,
						"inaccessible link, err: ", err, " url: ", href, "resp: ", resp)
					inaccessible = true
				}

				linkMu.Lock()
				if strings.HasPrefix(href, "http") {
					if strings.HasPrefix(href, baseURL) {
						link.InternalLinks++
					} else {
						link.ExternalLinks++
					}
				} else {
					link.InternalLinks++
				}
				if inaccessible {
					link.InaccessibleLinkCount++
					link.InaccessibleLink = append(link.InaccessibleLink, fullURL)
				}
				linkMu.Unlock()
			}
		}()
	}

	// select the element
	// then get the href
	// ignore the #
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if !exists || strings.HasPrefix(href, "#") {
			return
		}

		// add in to map to get the distinct links
		_, ok := distinctLinks[href]
		if !ok {
			distinctLinks[href] = nil
			linkJobs <- href
		}
	})

	close(linkJobs)
	wg.Wait()
	return link
}

func NewAnalyser(ctr container.Container, cfg bootstrap.Config) Analyser {
	return &analyser{
		ctr:    ctr,
		config: cfg,
	}
}
func normalizeURL(url string) string {
	return strings.TrimSuffix(url, "/")
}

func resolveURL(base, href string) string {
	if strings.HasPrefix(href, "http") {
		return href
	}
	if strings.HasPrefix(href, "/") {
		return base + href
	}
	return base + "/" + href
}
