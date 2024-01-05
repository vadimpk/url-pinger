package service

import "net/url"

type urlValidator struct {
}

func NewURLValidator() URLValidator {
	return &urlValidator{}
}

func (u *urlValidator) ValidateURL(uri string) bool {
	_, err := url.ParseRequestURI(uri)
	return err == nil
}
