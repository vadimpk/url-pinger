package service

type urlValidator struct {
}

func NewURLValidator() URLValidator {
	return &urlValidator{}
}

func (u *urlValidator) ValidateURL(url string) (bool, error) {
	return true, nil
}
