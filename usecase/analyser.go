package usecase

import (
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/web-page-analysis/container"
	"github.com/web-page-analysis/domain"
	"strings"
	"sync"
)

const (
	title = "title"
)

type Analyser interface {
	CheckHtmlVersion(ctx context.Context, rawHtml string) string
	GetTitle(ctx context.Context, doc *goquery.Document) string
	CountHeading(ctx context.Context, doc *goquery.Document) map[string]int
	CountLinks(ctx context.Context, doc *goquery.Document, url string) domain.Link
	CheckAnyLogin(ctx context.Context, doc *goquery.Document) bool
}

type analyser struct {
	ctr container.Container
}

func (a analyser) CheckHtmlVersion(ctx context.Context, rawHTML string) string {
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
	var isExist bool
	// Check for login form
	doc.Find("form").Each(func(i int, s *goquery.Selection) {
		if s.Find("input[type='password']").Length() > 0 ||
			s.Find("input[type='email']").Length() > 0 ||
			s.Find("input[name='username']").Length() > 0 ||
			s.Text() != "" && strings.Contains(strings.ToLower(s.Text()), "log in") {
			isExist = true
		}
	})
	return isExist
}

func (a analyser) GetTitle(ctx context.Context, doc *goquery.Document) string {
	return doc.Find(title).Text()
}

func (a analyser) CountHeading(ctx context.Context, doc *goquery.Document) map[string]int {
	headingsMap := map[string]int{}
	for i := 1; i <= 6; i++ {
		tag := fmt.Sprintf("h%d", i)
		headingsMap[tag] = doc.Find(tag).Length()
	}

	return headingsMap
}

func (a analyser) CountLinks(ctx context.Context, doc *goquery.Document, baseURL string) domain.Link {
	var (
		link        domain.Link
		linkMu      sync.Mutex
		linkJobs    = make(chan string)
		wg          sync.WaitGroup
		workerCount = 200
	)

	link.InaccessibleLink = make([]string, 0)

	baseURL = normalizeURL(baseURL)

	// Start worker pool
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for href := range linkJobs {
				var inaccessible bool
				fullURL := resolveURL(baseURL, href)

				// Check link accessibility
				resp, err := a.ctr.OBAdapter.Get(ctx, fullURL)
				if resp != nil && resp.Body != nil {
					defer resp.Body.Close()
				}

				if err != nil || resp != nil && (resp.StatusCode > 300 || resp.StatusCode < 200) {
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

	// Collect jobs
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if !exists || strings.HasPrefix(href, "#") {
			return
		}
		linkJobs <- href
	})

	close(linkJobs)
	wg.Wait()
	return link
}

func NewAnalyser(ctr container.Container) Analyser {
	return &analyser{
		ctr: ctr,
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
