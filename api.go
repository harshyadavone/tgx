package tgx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
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
			Message: "failed to read response body",
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

func makeMultipartReq(token, method string, params map[string]interface{}, paramName, path string) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/%s", token, method)

	fmt.Println("FilePath: ", path)
	fmt.Println("FileName: ", filepath.Base(path))
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, filepath.Base(path))
	if err != nil {
		return fmt.Errorf("failed to create form file: %w", err)
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return fmt.Errorf("failed to write file to form: %w", err)
	}

	for key, val := range params {
		strVal := fmt.Sprintf("%v", val)
		err = writer.WriteField(key, strVal)
		if err != nil {
			return fmt.Errorf("failed to write form field %q: %w", key, err)
		}
	}

	err = writer.Close()
	if err != nil {
		return fmt.Errorf("failed to close writer: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	fmt.Println(string(respBody))

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("unexpected response status: %d, body: %s", resp.StatusCode, string(respBody))
	}

	return nil
}

func IsAPIError(err error, errCode int) bool {
	if botErr, ok := err.(*BotError); ok {
		if apiErr, ok := botErr.Err.(*APIError); ok {
			return apiErr.Code == errCode
		}
	}
	return false
}
