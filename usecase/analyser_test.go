package usecase

import (
	"context"
	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"
	"github.com/web-page-analysis/bootstrap"
	"github.com/web-page-analysis/container"
	"strings"
	"testing"
)

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
	conf := bootstrap.Config{
		OutboundConf: bootstrap.OutboundConfig{
			DialTimeout:   60,
			RemoteTimeout: 60,
		},
	}
	ctr := container.Resolver(ctx, conf)
	analyser := NewAnalyser(*ctr)

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
	conf := bootstrap.Config{
		OutboundConf: bootstrap.OutboundConfig{
			DialTimeout:   60,
			RemoteTimeout: 60,
		},
	}
	ctr := container.Resolver(ctx, conf)
	analyser := NewAnalyser(*ctr)

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
	conf := bootstrap.Config{
		OutboundConf: bootstrap.OutboundConfig{
			DialTimeout:   60,
			RemoteTimeout: 60,
		},
	}
	ctr := container.Resolver(ctx, conf)
	analyser := NewAnalyser(*ctr)

	actual := analyser.CountHeading(ctx, docFromHTML(t, htmlForTitleCheck))
	assert.Equal(t, 2, actual["h1"])
	assert.Equal(t, 1, actual["h2"])
	assert.Equal(t, 1, actual["h3"])
	assert.Equal(t, 0, actual["h4"])
}

func TestCountLinks(t *testing.T) {
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
	conf := bootstrap.Config{
		OutboundConf: bootstrap.OutboundConfig{
			DialTimeout:   2000,
			RemoteTimeout: 2000,
		},
	}
	ctr := container.Resolver(ctx, conf)
	analyser := NewAnalyser(*ctr)

	actual := analyser.CountLinks(ctx, docFromHTML(t, htmlForTitleCheck), "http://abc.com/")
	assert.Equal(t, 3, actual.InternalLinks)
	assert.Equal(t, 3, actual.ExternalLinks)
	assert.Equal(t, 3, actual.InaccessibleLinkCount)
}
