package mta

import "fmt"

var (
	ErrAPIKeyRequired      = fmt.Errorf("API key required")
	ErrAPIKeyNotAuthorized = fmt.Errorf("API key not authorized")
	ErrClientRequired      = fmt.Errorf("client required")
)
