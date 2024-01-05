package entity

type URLStatus string

const (
	// URLStatusOK is the status for a URL that is reachable.
	URLStatusOK URLStatus = "OK"
	// URLStatusNotFound is the status for a URL that is not found.
	URLStatusNotFound URLStatus = "NOT_FOUND"
	// URLStatusError is the status for a URL that has an error.
	URLStatusError URLStatus = "ERROR"
	// URLStatusTimeout is the status for a URL that has timed out.
	URLStatusTimeout URLStatus = "TIMEOUT"
	// URLStatusFailed is the status for a URL that we failed to ping.
	URLStatusFailed URLStatus = "FAILED"
	// URLStatusUnknown is the status for a URL that we don't know the status of.
	URLStatusUnknown URLStatus = "UNKNOWN"
	// URLStatusInvalid is the status for a URL that is invalid.
	URLStatusInvalid URLStatus = "INVALID"
)
