package usecase

import (
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/web-page-analysis/container"
	"github.com/web-page-analysis/domain"
	"strings"
)

const (
	title = "title"
)

type Analyser interface {
	GetTitle(ctx context.Context, doc *goquery.Document) string
	CountHeading(ctx context.Context, doc *goquery.Document) map[string]int
	CountLinks(ctx context.Context, doc *goquery.Document, url string) domain.Link
	CheckAnyLogin(ctx context.Context, doc *goquery.Document) bool
}

type analyser struct {
	ctr container.Container
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

func (a analyser) CountLinks(ctx context.Context, doc *goquery.Document, url string) domain.Link {
	// Links analysis
	var link domain.Link
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, isExist := s.Attr("href")
		if isExist {
			if strings.HasPrefix(href, "http") {
				if strings.HasPrefix(href, url) {
					link.InternalLinks++
				} else {
					link.ExternalLinks++
				}
			} else {
				link.InternalLinks++
				if strings.HasPrefix(href, "/") {
					href = normalizeURL(url) + href
				} else {
					href = url + href
				}
			}
			// Check the link status
			resp, err := a.ctr.OBAdapter.Get(ctx, href)
			if err == nil {
				defer resp.Body.Close()
			}
			if err != nil || resp != nil && (resp.StatusCode > 300 || resp.StatusCode < 200) {
				link.InaccessibleLinkCount++
				link.InaccessibleLink = append(link.InaccessibleLink, href)
			}
		}

	})
	return link
}

func NewAnalyser(ctr container.Container) Analyser {
	return &analyser{
		ctr: ctr,
	}
}

func normalizeURL(u string) string {
	return strings.TrimSuffix(u, "/")
}
