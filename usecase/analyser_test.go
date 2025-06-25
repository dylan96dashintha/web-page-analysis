package usecase

import (
	"context"
	"errors"
	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"
	"github.com/web-page-analysis/bootstrap"
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
	conf := bootstrap.Config{}
	analyser := NewAnalyser(ctr, conf)

	actual := analyser.GetTitle(ctx, docFromHTML(t, htmlForTitleCheck))
	assert.Equal(t, "Test Page", actual)
}

func TestGetWithoutTitle(t *testing.T) {
	var (
		htmlForTitleCheckWithoutTitle = `
<!DOCTYPE html>
<html>
<body><h1>Hello World</h1></body>
</html>
`
	)
	ctx := context.Background()
	ctr := container.Container{OBAdapter: mockOutBoundConnection{}}
	conf := bootstrap.Config{
		AppConfig: bootstrap.AppConfig{
			WorkerCount: 200,
		},
	}
	analyser := NewAnalyser(ctr, conf)

	actual := analyser.GetTitle(ctx, docFromHTML(t, htmlForTitleCheckWithoutTitle))
	assert.Equal(t, "", actual)
}

func TestCountHeading(t *testing.T) {
	var (
		htmlForCountHeading = `
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
	conf := bootstrap.Config{
		AppConfig: bootstrap.AppConfig{
			WorkerCount: 200,
		},
	}
	analyser := NewAnalyser(ctr, conf)

	actual := analyser.CountHeading(ctx, docFromHTML(t, htmlForCountHeading))
	assert.Equal(t, 2, actual["h1"])
	assert.Equal(t, 1, actual["h2"])
	assert.Equal(t, 1, actual["h3"])
	assert.Equal(t, 0, actual["h4"])
}

func TestCountLinksWithAccessible(t *testing.T) {
	var (
		htmlForCountHeadingWithAccessible = `
<!DOCTYPE html>
<html>
<head>
    <title>Link Test Page</title>
</head>
<body>
    <h1>Welcome</h1>

    // internal links
    <a href="/about">About Us</a>
    <a href="/contact">Contact</a>
    <a href="privacy.html">Privacy Policy</a>

    // external links
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
	conf := bootstrap.Config{
		AppConfig: bootstrap.AppConfig{
			WorkerCount: 200,
		},
	}
	analyser := NewAnalyser(ctr, conf)

	actual := analyser.CountLinks(ctx, docFromHTML(t, htmlForCountHeadingWithAccessible), "http://abc.com/")
	assert.Equal(t, 3, actual.InternalLinks)
	assert.Equal(t, 3, actual.ExternalLinks)
	assert.Equal(t, 0, actual.InaccessibleLinkCount)
}

func TestCountLinksWithDuplicates(t *testing.T) {
	var (
		htmlForCountHeadingWithDuplicates = `
<!DOCTYPE html>
<html>
<head>
    <title>Link Test Page</title>
</head>
<body>
    <h1>Welcome</h1>

	// internal links
    <a href="/about">About Us</a>
    <a href="/contact">Contact</a>
    <a href="/contact">Privacy Policy</a>
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
	conf := bootstrap.Config{
		AppConfig: bootstrap.AppConfig{
			WorkerCount: 200,
		},
	}
	analyser := NewAnalyser(ctr, conf)

	actual := analyser.CountLinks(ctx, docFromHTML(t, htmlForCountHeadingWithDuplicates), "http://abc.com/")
	assert.Equal(t, 2, actual.InternalLinks)
	assert.Equal(t, 0, actual.ExternalLinks)
	assert.Equal(t, 0, actual.InaccessibleLinkCount)
}

func TestCountLinksWithHashLink(t *testing.T) {
	var (
		htmlForCountHeadingWithHashLink = `
<!DOCTYPE html>
<html>
<head>
    <title>Link Test Page</title>
</head>
<body>
    <h1>Welcome</h1>

	// internal links
    <a href="/about">About Us</a>
    <a href="/contact">Contact</a>
    <a href="#blog">Privacy Policy</a>
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
	conf := bootstrap.Config{
		AppConfig: bootstrap.AppConfig{
			WorkerCount: 200,
		},
	}
	analyser := NewAnalyser(ctr, conf)

	actual := analyser.CountLinks(ctx, docFromHTML(t, htmlForCountHeadingWithHashLink), "http://abc.com/")
	assert.Equal(t, 2, actual.InternalLinks)
	assert.Equal(t, 0, actual.ExternalLinks)
	assert.Equal(t, 0, actual.InaccessibleLinkCount)
}

func TestCountLinksWithInaccessible(t *testing.T) {
	var (
		htmlForCountHeadingWithInaccessible = `
<!DOCTYPE html>
<html>
<head>
    <title>Link Test Page</title>
</head>
<body>
    <h1>Welcome</h1>

    // internal links
    <a href="/about">About Us</a>
    <a href="/contact">Contact</a>
    <a href="privacy.html">Privacy Policy</a>

    // external links
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
	conf := bootstrap.Config{
		AppConfig: bootstrap.AppConfig{
			WorkerCount: 200,
		},
	}
	analyser := NewAnalyser(ctr, conf)

	actual := analyser.CountLinks(ctx, docFromHTML(t, htmlForCountHeadingWithInaccessible), "http://abc.com/")
	assert.Equal(t, 3, actual.InternalLinks)
	assert.Equal(t, 3, actual.ExternalLinks)
	assert.Equal(t, 6, actual.InaccessibleLinkCount)
}

func TestCountLinksWithInaccessibleByReturningError(t *testing.T) {
	var (
		htmlForCountHeadingWithInaccessibleByReturningError = `
<!DOCTYPE html>
<html>
<head>
    <title>Link Test Page</title>
</head>
<body>
    <h1>Welcome</h1>

    // internal links
    <a href="/about">About Us</a>
    <a href="/contact">Contact</a>
    <a href="privacy.html">Privacy Policy</a>

    // external links
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
	conf := bootstrap.Config{
		AppConfig: bootstrap.AppConfig{
			WorkerCount: 200,
		},
	}
	analyser := NewAnalyser(ctr, conf)

	actual := analyser.CountLinks(ctx, docFromHTML(t, htmlForCountHeadingWithInaccessibleByReturningError), "http://abc.com/")
	assert.Equal(t, 3, actual.InternalLinks)
	assert.Equal(t, 3, actual.ExternalLinks)
	assert.Equal(t, 6, actual.InaccessibleLinkCount)
}

func TestCheckHtmlVersion(t *testing.T) {
	var (
		htmlForHtmlVersion = `
<!DOCTYPE html>
<html>
<head>
    <title>Link Test Page</title>
</head>
<body>
    <h1>Welcome</h1>
</body>
</html>
`
	)
	ctx := context.Background()
	ctr := container.Container{OBAdapter: mockOutBoundConnection{}}
	conf := bootstrap.Config{
		AppConfig: bootstrap.AppConfig{
			WorkerCount: 200,
		},
	}
	analyser := NewAnalyser(ctr, conf)

	actual := analyser.CheckHtmlVersion(ctx, htmlForHtmlVersion)
	assert.Equal(t, "HTML5", actual)
}

func TestCheckAnyLogin(t *testing.T) {
	var (
		htmlForLoginCheck = `
<!DOCTYPE html>
<html>
<head>
    <title>Login Page</title>
</head>
<body>
    <h1>Please Login</h1>

    <form action="/login" method="post">
        <label for="username">Username:</label>
        <input type="text" id="username" name="username" />
    </form>

    <a href="/forgot-password">Forgot Password?</a>
    <a href="/signup">Sign up</a>
</body>
</html>
`
	)
	ctx := context.Background()
	ctr := container.Container{OBAdapter: mockOutBoundConnection{}}
	conf := bootstrap.Config{
		AppConfig: bootstrap.AppConfig{
			WorkerCount: 200,
		},
	}
	analyser := NewAnalyser(ctr, conf)

	actual := analyser.CheckAnyLogin(ctx, docFromHTML(t, htmlForLoginCheck))
	assert.Equal(t, true, actual)
}

func TestCheckAnyLoginWithPassword(t *testing.T) {
	var (
		htmlForLoginCheckWithPassword = `
<!DOCTYPE html>
<html>
<head>
    <title>Login Page</title>
</head>
<body>
    <h1>Please Login</h1>

    <form action="/login" method="post">
         <input type="Password" id="password" name="password" />
    </form>

    <a href="/forgot-password">Forgot Password?</a>
    <a href="/signup">Sign up</a>
</body>
</html>
`
	)
	ctx := context.Background()
	ctr := container.Container{OBAdapter: mockOutBoundConnection{}}
	conf := bootstrap.Config{
		AppConfig: bootstrap.AppConfig{
			WorkerCount: 200,
		},
	}
	analyser := NewAnalyser(ctr, conf)

	actual := analyser.CheckAnyLogin(ctx, docFromHTML(t, htmlForLoginCheckWithPassword))
	assert.Equal(t, true, actual)
}
