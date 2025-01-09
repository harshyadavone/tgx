package tgx

import (
	"github.com/harshyadavone/tgx/models"
)

func (b *Bot) OnMessage(handler handler) {
	b.messageHandler = handler
}

func (ctx *Context) Reply(req *SendMessageRequest) error {
	payload := map[string]interface{}{
		"chat_id": ctx.ChatID,
		"text":    req.Text,
	}

	if req.ParseMode != "" {
		payload["parse_mode"] = req.ParseMode
	}

	if req.ReplyMarkup != nil {
		payload["reply_markup"] = req.ReplyMarkup
	}

	return ctx.makeRequest("sendMessage", payload)
}

func (ctx *Context) ReplyWithInlineKeyboard(text string, buttons [][]models.InlineKeyboardButton) error {
	return ctx.makeRequest("sendMessage", map[string]any{
		"chat_id": ctx.ChatID,
		"text":    text,
		"reply_markup": map[string]any{
			"inline_keyboard": buttons,
		},
	})
}
