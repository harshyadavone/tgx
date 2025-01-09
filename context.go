package tgx

import "github.com/harshyadavone/tgx/models"

type Context struct {
	Text     string
	UserID   int64
	Username string
	ChatID   int64
	bot      *Bot
	Args     []string
}

type CallbackContext struct {
	QueryID  string
	Data     string
	Message  *models.Message
	UserID   int64
	Username string
	bot      *Bot
}
