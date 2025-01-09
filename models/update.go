package models

type Update struct {
	UpdateId      int            `json:"update_id"`
	Message       *Message       `json:"message"`
	CallbackQuery *CallbackQuery `json:"callback_query"`
	EditedMessage *Message       `json:"edited_message"`
	InlineQuery   *InlineQuery   `json:"inline_query"`
}

type InlineQuery struct {
	Id       int64  `json:"id"`
	From     User   `json:"from"`
	Query    string `json:"query"`
	Offset   string `json:"offset"`
	ChatType string `json:"chat_type"`
}

type CallbackQuery struct {
	ID      string   `json:"id"`
	From    User     `json:"from"`
	Message *Message `json:"message"`
	Data    string   `json:"data"`
}

type Message struct {
	MessageId      int64                 `json:"message_id"`
	From           User                  `json:"from"`
	Chat           Chat                  `json:"chat"`
	Text           string                `json:"text"`
	ReplyToMessage *Message              `json:"reply_to_message"`
	ReplyMarkup    *InlineKeyboardMarkup `json:"reply_markup"`
	Animation      *Animation            `json:"animation"`
	Audio          *Audio                `json:"audio"`
	Document       *Document             `json:"document"`
	Photo          []*PhotoSize          `json:"photo"`
	Sticker        *Sticker              `json:"sticker"`
	Video          *Video                `json:"video"`
	VideoNote      *VideoNote            `json:"video_note"`
	Voice          *Voice                `json:"voice"`
	Caption        string                `json:"caption"`
}
