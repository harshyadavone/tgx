package tgx

import (
	"fmt"

	"github.com/harshyadavone/tgx/models"
)

func (b *Bot) OnCallback(data string, handler callbackHandler) {
	b.callbackHandlers[data] = handler
}

func (b *Bot) handleCallbackQuery(cb *models.CallbackQuery) error {
	ctx := &CallbackContext{
		QueryID:  cb.ID,
		Data:     cb.Data,
		Message:  cb.Message,
		UserID:   cb.From.Id,
		Username: cb.From.Username,
		bot:      b,
	}

	if handler, ok := b.callbackHandlers[cb.Data]; ok {
		fmt.Println("Callback handler called")
		handler(ctx)
	}

	return nil
}

func (ctx *CallbackContext) Answer(text string) error {
	return ctx.makeRequest("sendMessage", map[string]any{
		"callback_query_id": ctx.QueryID,
		"text":              text,
		"chat_id":           ctx.Message.Chat.Id,
	})
}

func (ctx *CallbackContext) EditMessage(newText string) error {
	return ctx.makeRequest("editMessageText", map[string]any{
		"chat_id":    ctx.Message.Chat.Id,
		"message_id": ctx.Message.MessageId,
		"text":       newText,
	})
}
