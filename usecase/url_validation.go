package usecase

import (
	"context"
	log "github.com/sirupsen/logrus"
	"net/url"
)

const (
	validationPrefix = "usecase.url_validation "
)

type Validation interface {
	IsValidUrl(ctx context.Context, urlString string) bool
}

type validation struct{}

func (v validation) IsValidUrl(ctx context.Context, urlString string) bool {
	log.WithContext(ctx).Info(validationPrefix, "start to validate the url", urlString)
	parsedURL, err := url.ParseRequestURI(urlString)
	if err != nil {
		return false
	}
	if parsedURL != nil && (parsedURL.Scheme == "" || parsedURL.Host == "") {
		return false
	}
	return true
}

func NewValidation() Validation {
	return &validation{}
}
