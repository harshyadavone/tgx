package tgx

import "fmt"

type BotError struct {
	Code    int
	Message string
	Err     error
}

func (e *BotError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

type APIError struct {
	Code        int    `json:"error_code"`
	Description string `json:"description"`
	Parameters  struct {
		RetryAfter int `json:"retry_after,omitempty"`
	} `json:"parameters,omitempty"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("Telegram API error (code: %d): %s", e.Code, e.Description)
}
