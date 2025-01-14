package tgx

import "github.com/harshyadavone/tgx/models"

type Context struct {
	Text      string
	Photo     []*models.PhotoSize
	Video     *models.Video
	Voice     *models.Voice
	Document  *models.Document
	Sticker   *models.Sticker
	Animation *models.Animation
	Audio     *models.Audio
	VideoNote *models.VideoNote
	Args      []string
	UserID    int64
	Username  string
	MessageId int64
	ChatID    int64
	bot       *Bot
}

type CallbackContext struct {
	QueryID  string
	Data     string
	Message  *models.Message
	UserID   int64
	Username string
	bot      *Bot
}
