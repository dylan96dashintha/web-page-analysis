package usecase

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsValidUrl(t *testing.T) {
	var (
		urlString = "https://www.google.com"
	)
	ctx := context.Background()
	urlValidator := NewValidation()

	actual := urlValidator.IsValidUrl(ctx, urlString)
	assert.Equal(t, true, actual)
}

func TestIsValidUrlError(t *testing.T) {
	var (
		urlString = "htpp//google.com"
	)
	ctx := context.Background()
	urlValidator := NewValidation()

	actual := urlValidator.IsValidUrl(ctx, urlString)
	assert.Equal(t, false, actual)
}

func TestIsValidUrlNoHost(t *testing.T) {
	var (
		urlString = "http://"
	)
	ctx := context.Background()
	urlValidator := NewValidation()

	actual := urlValidator.IsValidUrl(ctx, urlString)
	assert.Equal(t, false, actual)
}

func TestIsValidUrlNoSchema(t *testing.T) {
	var (
		urlString = "://example.com"
	)
	ctx := context.Background()
	urlValidator := NewValidation()

	actual := urlValidator.IsValidUrl(ctx, urlString)
	assert.Equal(t, false, actual)
}
