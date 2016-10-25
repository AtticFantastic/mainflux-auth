package domain

import "fmt"

// ServiceError is used as a wrapper for all types of errors originated
// throughout the service. Code will be used as API response status.
type ServiceError struct {
	Code    int
	Message string
}

func (e *ServiceError) Error() string {
	if e == nil {
		return "<nil>"
	}

	return fmt.Sprintf("code: %d, message: %s", e.Code, e.Message)
}
