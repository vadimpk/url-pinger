package entity

type URLStatus string

const (
	URLStatusOK       URLStatus = "OK"
	URLStatusNotFound URLStatus = "NOT_FOUND"
	URLStatusError    URLStatus = "ERROR"
	URLStatusUnknown  URLStatus = "UNKNOWN"
)
