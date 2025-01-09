package tgx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

var client = &http.Client{
	Timeout: 30 * time.Second,
}

type TelegramResponse struct {
	Ok          bool            `json:"ok"`
	Result      json.RawMessage `json:"result"`
	Description string          `json:"description"`
	ErrorCode   int             `json:"error_code"`
}

func (ctx *Context) makeRequest(method string, params map[string]interface{}) error {
	return makeAPIRequest(ctx.bot.token, method, params)
}
func (ctx *CallbackContext) makeRequest(method string, params map[string]interface{}) error {
	return makeAPIRequest(ctx.bot.token, method, params)
}

func makeAPIRequest(token, method string, params map[string]interface{}) error {
	_, err := makeAPIRequestWithResult(token, method, params)
	return err
}

func makeAPIRequestWithResult(token, method string, params map[string]interface{}) (json.RawMessage, error) {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/%s", token, method)

	body, err := json.Marshal(params)
	if err != nil {
		return nil, &BotError{
			Code:    http.StatusInternalServerError,
			Message: "Failed to encode request parameters",
			Err:     err,
		}
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, &BotError{
			Code:    http.StatusInternalServerError,
			Message: "Failed to create request",
			Err:     err,
		}
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, &BotError{
			Code:    http.StatusServiceUnavailable,
			Message: "Failed to send request",
			Err:     err,
		}
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &BotError{
			Code:    http.StatusServiceUnavailable,
			Message: "Failed to send request",
			Err:     err,
		}
	}

	var telegramResp TelegramResponse
	if err = json.Unmarshal(respBody, &telegramResp); err != nil {
		return nil, &BotError{
			Code:    http.StatusInternalServerError,
			Message: "Failed to parse response",
			Err:     err,
		}
	}

	if !telegramResp.Ok {
		apiError := &APIError{
			Code:        telegramResp.ErrorCode,
			Description: telegramResp.Description,
		}

		switch apiError.Code {
		case 429:
			return nil, &BotError{
				Code:    http.StatusTooManyRequests,
				Message: fmt.Sprintf("Rate limited. Retry after %d seconds", apiError.Parameters.RetryAfter),
				Err:     apiError,
			}
		case 400:
			return nil, &BotError{
				Code:    http.StatusBadRequest,
				Message: "Invalid request to Telegram API",
				Err:     apiError,
			}
		case 401:
			return nil, &BotError{
				Code:    http.StatusUnauthorized,
				Message: "Invalid bot token",
				Err:     apiError,
			}
		case 403:
			return nil, &BotError{
				Code:    http.StatusForbidden,
				Message: "Bot lacks necessary permissions",
				Err:     apiError,
			}
		default:
			return nil, &BotError{
				Code:    resp.StatusCode,
				Message: "Telegram API error",
				Err:     apiError,
			}
		}
	}

	return telegramResp.Result, nil
}

func IsAPIError(err error, errCode int) bool {
	if botErr, ok := err.(*BotError); ok {
		if apiErr, ok := botErr.Err.(*APIError); ok {
			return apiErr.Code == errCode
		}
	}
	return false
}
