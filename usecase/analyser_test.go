package usecase

import (
	"context"
	"errors"
	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"
	"github.com/web-page-analysis/container"
	"io"
	"net/http"
	"strings"
	"testing"
)

type mockResolver struct{}
type mockOutBoundConnection struct {
}

var (
	mockOutboundResp  *http.Response
	mockOutBoundError error
)

func (o mockOutBoundConnection) Get(ctx context.Context, url string) (*http.Response, error) {
	return mockOutboundResp, mockOutBoundError

}

func docFromHTML(t *testing.T, html string) *goquery.Document {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		t.Fatalf("Failed to parse HTML: %v", err)
	}
	return doc
}

func TestGetTitle(t *testing.T) {
	var (
		htmlForTitleCheck = `
<!DOCTYPE html>
<html>
<head><title>Test Page</title></head>
<body><h1>Hello World</h1></body>
</html>
`
	)
	ctx := context.Background()
	ctr := container.Container{OBAdapter: mockOutBoundConnection{}}
	analyser := NewAnalyser(ctr)

	actual := analyser.GetTitle(ctx, docFromHTML(t, htmlForTitleCheck))
	assert.Equal(t, "Test Page", actual)
}

func TestGetWithoutTitle(t *testing.T) {
	var (
		htmlForTitleCheck = `
<!DOCTYPE html>
<html>
<body><h1>Hello World</h1></body>
</html>
`
	)
	ctx := context.Background()
	ctr := container.Container{OBAdapter: mockOutBoundConnection{}}
	analyser := NewAnalyser(ctr)

	actual := analyser.GetTitle(ctx, docFromHTML(t, htmlForTitleCheck))
	assert.Equal(t, "", actual)
}

func TestCountHeading(t *testing.T) {
	var (
		htmlForTitleCheck = `
<!DOCTYPE html>
<html>
<head><title>Test Page</title></head>
<body><h1>Hello World</h1></body>
<body><h1>Hello World</h1></body>
<body><h2>Hello World</h1></body>
<body><h3>Hello World</h1></body>
</html>
`
	)
	ctx := context.Background()
	ctr := container.Container{OBAdapter: mockOutBoundConnection{}}
	analyser := NewAnalyser(ctr)

	actual := analyser.CountHeading(ctx, docFromHTML(t, htmlForTitleCheck))
	assert.Equal(t, 2, actual["h1"])
	assert.Equal(t, 1, actual["h2"])
	assert.Equal(t, 1, actual["h3"])
	assert.Equal(t, 0, actual["h4"])
}

func TestCountLinksWithAccessible(t *testing.T) {
	var (
		htmlForTitleCheck = `
<!DOCTYPE html>
<html>
<head>
    <title>Link Test Page</title>
</head>
<body>
    <h1>Welcome</h1>

    <!-- Internal links -->
    <a href="/about">About Us</a>
    <a href="/contact">Contact</a>
    <a href="privacy.html">Privacy Policy</a>

    <!-- External links -->
    <a href="https://example.com">Example</a>
    <a href="http://external.org/page">External Page</a>
    <a href="https://www.google.com/search?q=go">Google Search</a>
</body>
</html>
`
	)
	ctx := context.Background()
	ctr := container.Container{OBAdapter: mockOutBoundConnection{}}
	mockOutboundResp = &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader("")),
	}
	analyser := NewAnalyser(ctr)

	actual := analyser.CountLinks(ctx, docFromHTML(t, htmlForTitleCheck), "http://abc.com/")
	assert.Equal(t, 3, actual.InternalLinks)
	assert.Equal(t, 3, actual.ExternalLinks)
	assert.Equal(t, 0, actual.InaccessibleLinkCount)
}

func TestCountLinksWithInaccessible(t *testing.T) {
	var (
		htmlForTitleCheck = `
<!DOCTYPE html>
<html>
<head>
    <title>Link Test Page</title>
</head>
<body>
    <h1>Welcome</h1>

    <!-- Internal links -->
    <a href="/about">About Us</a>
    <a href="/contact">Contact</a>
    <a href="privacy.html">Privacy Policy</a>

    <!-- External links -->
    <a href="https://example.com">Example</a>
    <a href="http://external.org/page">External Page</a>
    <a href="https://www.google.com/search?q=go">Google Search</a>
</body>
</html>
`
	)
	ctx := context.Background()
	ctr := container.Container{OBAdapter: mockOutBoundConnection{}}
	mockOutboundResp = &http.Response{
		StatusCode: 404,
		Body:       io.NopCloser(strings.NewReader("")),
	}
	analyser := NewAnalyser(ctr)

	actual := analyser.CountLinks(ctx, docFromHTML(t, htmlForTitleCheck), "http://abc.com/")
	assert.Equal(t, 3, actual.InternalLinks)
	assert.Equal(t, 3, actual.ExternalLinks)
	assert.Equal(t, 6, actual.InaccessibleLinkCount)
}

func TestCountLinksWithInaccessibleByReturninigError(t *testing.T) {
	var (
		htmlForTitleCheck = `
<!DOCTYPE html>
<html>
<head>
    <title>Link Test Page</title>
</head>
<body>
    <h1>Welcome</h1>

    <!-- Internal links -->
    <a href="/about">About Us</a>
    <a href="/contact">Contact</a>
    <a href="privacy.html">Privacy Policy</a>

    <!-- External links -->
    <a href="https://example.com">Example</a>
    <a href="http://external.org/page">External Page</a>
    <a href="https://www.google.com/search?q=go">Google Search</a>
</body>
</html>
`
	)
	ctx := context.Background()
	ctr := container.Container{OBAdapter: mockOutBoundConnection{}}
	mockOutBoundError = errors.New("mock outbound error")
	analyser := NewAnalyser(ctr)

	actual := analyser.CountLinks(ctx, docFromHTML(t, htmlForTitleCheck), "http://abc.com/")
	assert.Equal(t, 3, actual.InternalLinks)
	assert.Equal(t, 3, actual.ExternalLinks)
	assert.Equal(t, 6, actual.InaccessibleLinkCount)
}

func TestCheckHtmlVersion(t *testing.T) {
	var (
		htmlForTitleCheck = `
<!DOCTYPE html>
<html>
<head>
    <title>Link Test Page</title>
</head>
<body>
    <h1>Welcome</h1>

    <!-- Internal links -->
    <a href="/about">About Us</a>
    <a href="/contact">Contact</a>
    <a href="privacy.html">Privacy Policy</a>

    <!-- External links -->
    <a href="https://example.com">Example</a>
    <a href="http://external.org/page">External Page</a>
    <a href="https://www.google.com/search?q=go">Google Search</a>
</body>
</html>
`
	)
	ctx := context.Background()
	ctr := container.Container{OBAdapter: mockOutBoundConnection{}}
	analyser := NewAnalyser(ctr)

	actual := analyser.CheckHtmlVersion(ctx, htmlForTitleCheck)
	assert.Equal(t, "HTML5", actual)
}
