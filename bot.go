package tgx

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/harshyadavone/tgx/models"
	"github.com/harshyadavone/tgx/pkg/logger"
)

const (
	ParseModeMarkdown = "MarkdownV2"
	ParseModeHTML     = "HTML"
)

type ErrorHandler func(ctx *Context, err error)
type handler func(ctx *Context) error
type callbackHandler func(ctx *CallbackContext) error

type Bot struct {
	token      string
	webhookURL string

	messageHandler   handler
	commandHandler   map[string]handler
	callbackHandlers map[string]callbackHandler
	errorHandler     ErrorHandler

	logger  logger.Logger
	builder *ParamBuilder
}

type WebhookInfo struct {
	URL                          string   `json:"url"`
	HasCustomCertificate         bool     `json:"has_custom_certificate"`
	PendingUpdateCount           int      `json:"pending_update_count"`
	IPAddress                    string   `json:"ip_address,omitempty"`
	LastErrorDate                int64    `json:"last_error_date,omitempty"`
	LastErrorMessage             string   `json:"last_error_message,omitempty"`
	LastSynchronizationErrorDate int64    `json:"last_synchronization_error_date,omitempty"`
	MaxConnections               int      `json:"max_connections,omitempty"`
	AllowedUpdates               []string `json:"allowed_updates,omitempty"`
}

func NewBot(token, webhookURL string, logger logger.Logger) *Bot {
	return &Bot{
		token:            token,
		webhookURL:       webhookURL,
		commandHandler:   make(map[string]handler),
		callbackHandlers: make(map[string]callbackHandler),
		logger:           logger,
		errorHandler:     defaultErrorHandler,
		builder:          NewParamBuilder(),
	}
}

func defaultErrorHandler(ctx *Context, err error) {
	ctx.bot.logger.Error("Bot error:", err)
	payload := &SendMessageRequest{
		ChatId: ctx.ChatID,
		Text:   "Sorry, something went wrong. Please try again later.",
	}
	ctx.ReplyWithOpts(payload)
}

func (b *Bot) SetWebhook() error {
	return makeAPIRequest(b.token, "setWebhook", map[string]interface{}{
		"url": b.webhookURL,
	})
}

func (b *Bot) DeleteWebhook() error {
	return makeAPIRequest(b.token, "deleteWebhook", map[string]interface{}{})
}

func (b *Bot) GetWebhookInfo() (*WebhookInfo, error) {
	result, err := makeAPIRequestWithResult(b.token, "getWebhookInfo", nil)
	if err != nil {
		return nil, err
	}

	var webhookInfo WebhookInfo
	if err := json.Unmarshal(result, &webhookInfo); err != nil {
		return nil, &BotError{
			Code:    http.StatusBadRequest,
			Message: "failed to decode webhook info",
			Err:     err,
		}

	}
	return &webhookInfo, nil
}

func (b *Bot) GetMe() (*models.User, error) {
	result, err := makeAPIRequestWithResult(b.token, "getMe", nil)
	if err != nil {
		return nil, err
	}

	var user models.User
	if err := json.Unmarshal(result, &user); err != nil {
		return nil, &BotError{
			Code:    http.StatusBadRequest,
			Message: "failed to decode webhook info",
			Err:     err,
		}

	}
	return &user, nil
}

func (b *Bot) logOut() error {
	return makeAPIRequest(b.token, "getMe", nil)
}

func (b *Bot) close() error {
	return makeAPIRequest(b.token, "close", nil)
}

func (b *Bot) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		b.logger.Error("Invalid HTTP method: %s", r.Method)
		http.Error(w, "Only POST requests are allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		b.logger.Error("Code : %d,\nMessage: %s,\nErr: %v", http.StatusBadRequest, "Failed to read request body", err)
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var update models.Update
	if err := json.Unmarshal(body, &update); err != nil {
		b.logger.Error("Code : %d,\nMessage: %s,\nErr: %v", http.StatusBadRequest, "Failed to decode body", err)
		http.Error(w, "Failed to decode update", http.StatusBadRequest)
		return
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				b.logger.Error("Panic recovered in update handler: %v", r)
			}
		}()

		if update.Message != nil {
			if err := b.handleMessageUpdate(update.Message); err != nil {
				b.logger.Error("Error handling message update: %v", err)
			}
		} else if update.CallbackQuery != nil {
			if err := b.handleCallbackQuery(update.CallbackQuery); err != nil {
				b.logger.Error("Error handling callback query: %v", err)
			}
		} else {
			b.logger.Warn("Received update with no message or callback query")
		}
	}()

	w.WriteHeader(http.StatusOK)
}

func (b *Bot) handleMessageUpdate(message *models.Message) error {

	if message.Text == "" {
		return &BotError{
			Code:    http.StatusBadRequest,
			Message: "Empty message received",
		}
	}

	defer func() {
		if r := recover(); r != nil {
			b.logger.Error("Panic recovered in handleMessageUpdate: %v", r)
		}
	}()

	ctx := &Context{
		Text:      message.Text,
		UserID:    message.From.Id,
		Username:  message.From.Username,
		MessageId: message.MessageId,
		ChatID:    message.Chat.Id,
		bot:       b,
	}

	if strings.HasPrefix(message.Text, "/") {

		parts := strings.Split(message.Text, " ")
		if len(parts) < 1 {
			b.logger.Error("Not a valid command")
			return &BotError{
				Code:    http.StatusInternalServerError,
				Message: "Failed to parse command",
				Err:     fmt.Errorf("Failed to parse command"),
			}
		}

		command := strings.Split(parts[0], "/")
		if len(command) <= 1 {
			b.logger.Error("Not a valid command")
			return &BotError{
				Code:    http.StatusInternalServerError,
				Message: "Not a valid command",
				Err:     fmt.Errorf("Failed to split parts"),
			}
		}

		if len(parts) > 1 {
			args := parts[1:]
			ctx.Args = args
		}

		b.logger.Debug("Received message: %s", message.Text)
		b.logger.Debug("Parsed command: %s", command[1])
		if len(ctx.Args) > 0 {
			b.logger.Debug("Arguments: [%s]", strings.Join(ctx.Args, ", "))
		}

		if handler, ok := b.commandHandler[command[1]]; ok {
			b.logger.Info("Executing command: %s", command[1])

			err := b.safeExecute(ctx, handler)
			if err != nil {
				return err
			}
		} else {
			return &BotError{
				Code:    http.StatusNotFound,
				Message: "Unknown command",
				Err:     fmt.Errorf("command '%s' not found", command[1]),
			}
		}

	} else if b.messageHandler != nil {
		if err := b.safeExecute(ctx, b.messageHandler); err != nil {
			return err
		}
	}
	return nil
}

func (b *Bot) safeExecute(ctx *Context, handler handler) error {
	defer func() {
		if r := recover(); r != nil {
			b.logger.Error("Panic in handler execution:", r)
			err := fmt.Errorf("handler panic: %v", r)
			if b.errorHandler != nil {
				b.errorHandler(ctx, err)
			}
		}
	}()
	err := handler(ctx)
	if err != nil {
		switch {
		case IsAPIError(err, 403):
			b.logger.Warn("Bot blocked by user: %d", ctx.UserID)
			return err
		case IsAPIError(err, 429):
			b.logger.Info("Rate limited")
			return err
		default:
			if b.errorHandler != nil {
				b.errorHandler(ctx, err)
			}
			return err
		}
	}
	return nil
}

// SendMessage

func (b *Bot) SendMessage(chatID int64, text string) error {
	return makeAPIRequest(b.token, "sendMessage", map[string]interface{}{
		"chat_id": chatID,
		"text":    text,
	})
}

func (b *Bot) SendMessageWithOpts(req *SendMessageRequest) error {
	payload := map[string]interface{}{
		"chat_id": req.ChatId,
		"text":    req.Text,
	}

	if req.ParseMode != "" {
		if req.ParseMode != ParseModeHTML && req.ParseMode != ParseModeMarkdown {
			return &BotError{
				Code:    http.StatusBadRequest,
				Message: "Parse mode can be only 'MarkdownV2' or 'HTML'",
				Err:     fmt.Errorf("Parse mode can be only 'MarkdownV2' or 'HTML'"),
			}
		}
		payload["parse_mode"] = req.ParseMode
	}

	if req.ReplyMarkup != nil {
		payload["reply_markup"] = req.ReplyMarkup
	}

	if req.ReplyParams != nil && req.ReplyParams.MessageId != 0 {
		replyParam := map[string]interface{}{
			"message_id": req.ReplyParams.MessageId,
		}
		if req.ReplyParams.ChatId != 0 {
			replyParam["chat_id"] = req.ReplyParams.ChatId
		}
		payload["reply_parameters"] = replyParam
	}

	return makeAPIRequest(b.token, "sendMessage", payload)
}

// Forward Message

func (b *Bot) ForwardMessage(chatId, fromChatId, messageId int64) error {
	return makeAPIRequest(b.token, "forwardMessage", map[string]interface{}{
		"chat_id":      chatId,
		"from_chat_id": fromChatId,
		"message_id":   messageId,
	})
}

func (b *Bot) ForwardMessageWithOpts(req *ForwardMessageRequest) error {
	payload := map[string]interface{}{
		"chat_id":      req.ChatId,
		"from_chat_id": req.FromChatId,
		"message_id":   req.MessageId,
	}

	if req.DisableNotification {
		payload["disable_notification"] = req.DisableNotification
	}

	if req.ProtectContent {
		payload["protect_content"] = req.ProtectContent
	}

	return makeAPIRequest(b.token, "forwardMessage", payload)
}

// ForwardMessages

func (b *Bot) ForwardMessages(chatId, fromChatId int64, messageId []int64) error {
	return makeAPIRequest(b.token, "forwardMessages", map[string]interface{}{
		"chat_id":      chatId,
		"from_chat_id": fromChatId,
		"message_id":   messageId,
	})
}

func (b *Bot) ForwardMessagesWithOpts(req *ForwardMessagesRequest) error {
	payload := map[string]interface{}{
		"chat_id":      req.ChatId,
		"from_chat_id": req.FromChatId,
		"message_id":   req.MessageIds,
	}

	if req.DisableNotification {
		payload["disable_notification"] = req.DisableNotification
	}

	if req.ProtectContent {
		payload["protect_content"] = req.ProtectContent
	}

	return makeAPIRequest(b.token, "forwardMessages", payload)
}

// CopyMessage

func (b *Bot) CopyMessage(chatId, fromChatId, messageId int64) error {
	return makeAPIRequest(b.token, "copyMessage", map[string]interface{}{
		"chat_id":      chatId,
		"from_chat_id": fromChatId,
		"message_id":   messageId,
	})
}

func (b *Bot) CopyMessageWithOpts(req *CopyMessageRequest) error {
	payload := map[string]interface{}{
		"chat_id":      req.ChatId,
		"from_chat_id": req.FromChatId,
		"message_id":   req.MessageId,
	}

	if req.Caption != "" {
		payload["caption"] = req.Caption
	}

	if req.ParseMode != "" {
		if req.ParseMode != ParseModeHTML && req.ParseMode != ParseModeMarkdown {
			return &BotError{
				Code:    http.StatusBadRequest,
				Message: "Parse mode can be only 'MarkdownV2' or 'HTML'",
				Err:     fmt.Errorf("Parse mode can be only 'MarkdownV2' or 'HTML'"),
			}
		}
		payload["parse_mode"] = req.ParseMode
	}

	if req.ShowCaptionAboveMedia {
		payload["show_caption_above_media"] = req.ShowCaptionAboveMedia
	}

	if req.AllowPaidBroadCast {
		payload["allow_paid_broadcast"] = req.AllowPaidBroadCast
	}

	if req.DisableNotification {
		payload["disable_notification"] = req.DisableNotification
	}

	if req.ProtectContent {
		payload["protect_content"] = req.ProtectContent
	}

	if req.ReplyMarkup != nil {
		payload["reply_markup"] = req.ReplyMarkup
	}

	if req.ReplyParams != nil && req.ReplyParams.MessageId != 0 {
		replyParam := map[string]interface{}{
			"message_id": req.ReplyParams.MessageId,
		}
		if req.ReplyParams.ChatId != 0 {
			replyParam["chat_id"] = req.ReplyParams.ChatId
		}
		payload["reply_parameters"] = replyParam
	}

	return makeAPIRequest(b.token, "copyMessage", payload)
}

// CopyMessages

func (b *Bot) CopyMessages(chatId, fromChatId int64, messageId []int64) error {
	return makeAPIRequest(b.token, "copyMessages", map[string]interface{}{
		"chat_id":      chatId,
		"from_chat_id": fromChatId,
		"message_id":   messageId,
	})
}

func (b *Bot) CopyMessagesWithOpts(req *CopyMessagesRequest) error {
	payload := map[string]interface{}{
		"chat_id":      req.ChatId,
		"from_chat_id": req.FromChatId,
		"message_id":   req.MessageIds,
	}

	if req.DisableNotification {
		payload["disable_notification"] = req.DisableNotification
	}

	if req.ProtectContent {
		payload["protect_content"] = req.ProtectContent
	}

	if req.RemoveCaption {
		payload["remove_caption"] = req.RemoveCaption
	}

	return makeAPIRequest(b.token, "copyMessages", payload)
}

type sendOption func(url.Values)

var (
	OptCaption = func(caption string) sendOption {
		return func(v url.Values) {
			v.Set("caption", caption)
		}
	}

	OptReplyToMessageID = func(messageID int64) sendOption {
		return func(v url.Values) {
			v.Set("reply_to_message_id", strconv.FormatInt(messageID, 10))
		}
	}

	OptReplyMarkup = func(markup interface{}) sendOption {
		return func(v url.Values) {
			bytes, _ := json.Marshal(markup)
			v.Set("reply_markup", string(bytes))
		}
	}
)

// file_id or url
func (b *Bot) SendPhoto(req *SendPhotoRequest) error {
	builder := NewParamBuilder().
		Add("chat_id", req.ChatId).
		Add("caption", req.Caption).
		Add("disable_notification", req.DisableNotification).
		Add("reply_to_message_id", req.ReplyParams.MessageId)

	if req.ReplyMarkup != nil {
		replyBytes, _ := json.Marshal(req.ReplyMarkup)
		builder.Add("reply_markup", string(replyBytes))
	}

	params := builder.Build()
	return makeAPIRequest(b.token, "sendPhoto", params)
}

// file_path
func (b *Bot) SendPhotoFile(req *SendPhotoRequest) error {
	b.logger.Debug("Reached in sendPhotoFile function")
	b.builder.Add("chat_id", req.ChatId).
		Add("caption", req.Caption).
		Add("reply_to_message_id", req.ReplyParams.MessageId)

	b.logger.Debug("setting up reply_markup sendPhotoFile function")
	if req.ReplyMarkup != nil {
		replyBytes, _ := json.Marshal(req.ReplyMarkup)
		b.builder.Add("reply_markup", string(replyBytes))
	}

	b.logger.Debug("building params sendPhotoFile function")
	params := b.builder.Build()
	b.logger.Debug("sending request sendPhotoFile function")
	return makeMultipartReq(b.token, "sendPhoto", params, "photo", req.Photo)
}

// sendChatAction

func (b *Bot) SendChatAction(chatId int64, action string) error {
	return makeAPIRequest(b.token, "sendChatAction", map[string]interface{}{
		"chat_id": chatId,
		"action":  action,
	})
}
