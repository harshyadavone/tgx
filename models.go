package main

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

type Animation struct {
	FileId       string `json:"file_id"`
	FileUniqueId string `json:"file_unique_id"`
	Width        int64  `json:"width"`
	Height       int64  `json:"height"`
	Duration     int64  `json:"duration"`
	FileName     string `json:"file_name"`
	MimeType     string `json:"mime_type"`
	FileSize     int64  `json:"file_size"`
}

type Audio struct {
	FileId       string `json:"file_id"`
	FileUniqueId string `json:"file_unique_id"`
	Duration     int64  `json:"duration"`
	Performer    string `json:"performer"`
	Title        string `json:"title"`
	FileName     string `json:"file_name"`
	MimeType     int64  `json:"mime_type"`
	FileSize     int64  `json:"file_size"`
}

type Document struct {
	FileId       string `json:"file_id"`
	FileUniqueId string `json:"file_unique_id"`
	FileName     string `json:"file_name"`
	MimeType     int64  `json:"mime_type"`
	FileSize     int64  `json:"file_size"`
}

type PhotoSize struct {
	FileId       string `json:"file_id"`
	FileUniqueId string `json:"file_unique_id"`
	Width        int64  `json:"width"`
	Height       int64  `json:"height"`
	FileSize     int64  `json:"file_size"`
}

type Sticker struct {
	FileId       string     `json:"file_id"`
	FileUniqueId string     `json:"file_unique_id"`
	Type         string     `json:"type"` // regular, mask, custom_emoji
	Width        int64      `json:"width"`
	Height       int64      `json:"height"`
	IsAnimated   bool       `json:"is_animated"`
	IsVideo      bool       `json:"is_video"`
	Thumbnail    *PhotoSize `json:"thumbnail"`
	Emoji        string     `json:"emoji"`
	SetName      string     `json:"set_name"`
	FileSize     int64      `json:"file_size"`
}

type Video struct {
	FileId       string `json:"file_id"`
	FileUniqueId string `json:"file_unique_id"`
	Width        int64  `json:"width"`
	Height       int64  `json:"height"`
	Duration     int64  `json:"duration"`
	FileName     string `json:"file_name"`
	MimeType     int64  `json:"mime_type"`
	FileSize     int64  `json:"file_size"`
}

type VideoNote struct {
	FileId       string `json:"file_id"`
	FileUniqueId string `json:"file_unique_id"`
	Length       int64  `json:"length"`
	Duration     int64  `json:"duration"`
	FileSize     int64  `json:"file_size"`
}

type Voice struct {
	FileId       string `json:"file_id"`
	FileUniqueId string `json:"file_unique_id"`
	Duration     int64  `json:"duration"`
	MimeType     int64  `json:"mime_type"`
	FileSize     int64  `json:"file_size"`
}

type User struct {
	Id        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
}

type Chat struct {
	Id   int64  `json:"id"`
	Type string `json:"type"` // private, group, supergroup, channel
}

type Context struct {
	Text     string
	UserID   int64
	Username string
	ChatID   int64
	bot      *Bot
}

type CallbackContext struct {
	QueryID  string
	Data     string
	Message  *Message
	UserID   int64
	Username string
	bot      *Bot
}

type InlineKeyboardMarkup struct {
	InlineKeyboard [][]InlineKeyboardButton `json:"inline_keyboard"`
}

type InlineKeyboardButton struct {
	Text         string `json:"text"`
	URL          string `json:"url"` // tg:// URL to be opened when the button is pressed for ex. tg://user?id=<user_id>
	CallbackData string `json:"callback_data"`
}

type ReplyKeyboardMarkup struct {
	Keyboard              [][]KeyboardButton `json:"keyboard"`
	IsPersistent          bool               `json:"is_persistent"`
	ResizeKeyboard        bool               `json:"resize_keyboard,omitempty"`
	OneTimeKeyboard       bool               `json:"one_time_keyboard"`
	InputFieldPlaceholder string             `json:"input_field_placeholder"`
	Selective             bool               `json:"selective"`
}

type KeyboardButton struct {
	Text string `json:"text"`
}

type ReplyKeyboardRemove struct {
	RemoveKeyboard bool `json:"remove_keyboard"`
	Selective      bool `json:"selective"`
}

type ForceReply struct {
	ForceReply            bool   `json:"force_reply"`
	InputFieldPlaceholder string `json:"input_field_placeholder"`
	Selective             bool   `json:"selective"`
}

type Bot struct {
	token      string
	webhookURL string

	messageHandler   func(ctx *Context)
	commandHandler   map[string]func(ctx *Context)
	callbackHandlers map[string]func(ctx *CallbackContext)
}

type SendMessageRequest struct {
	ChatId      int64       `json:"chat_id"`
	Text        string      `json:"text"`
	ParseMode   string      `json:"parse_mode,omitempty"` // MarkdownV2 || HTML
	ReplyMarkup ReplyMarkup `json:"reply_markup,omitempty"`
}

type ReplyMarkup interface{}

type ForwardMessageRequest struct {
	ChatId     int64 `json:"chat_id"`
	FromChatId int64 `json:"from_chat_id"`
	MessageId  int64 `json:"message_id"`
}

type SendPhotoRequest struct {
	ChatId      int64       `json:"chat_id"`
	Photo       string      `json:"photo"` // pass a file_id of photo which is already uploaded on telegram server
	Caption     string      `json:"caption,omitempty"`
	ParseMode   string      `json:"parse_mode,omitempty"` // for caption
	HasSpoiler  bool        `json:"has_spoiler,omitempty"`
	ReplyMarkup ReplyMarkup `json:"reply_markup,omitempty"`
}

type SendAudioRequest struct {
	ChatId      int64       `json:"chat_id"`
	Audio       string      `json:"audio"`
	Caption     string      `json:"caption,omitempty"`
	Duration    int64       `json:"duration,omitempty"` // duration of audio i seconds
	Performer   string      `json:"performer,omitempty"`
	Title       string      `json:"title,omitempty"`
	ParseMode   string      `json:"parse_mode,omitempty"` // for caption
	ReplyMarkup ReplyMarkup `json:"reply_markup,omitempty"`
}

type SendVideoRequest struct {
	ChatId            int64       `json:"chat_id"`
	Video             string      `json:"video"`
	Caption           string      `json:"caption,omitempty"`
	Duration          int64       `json:"duration,omitempty"` // duration of video in seconds
	Width             int64       `json:"width"`
	Height            int64       `json:"height"`
	ParseMode         string      `json:"parse_mode,omitempty"` // for caption
	HasSpoiler        bool        `json:"has_spoiler,omitempty"`
	ReplyMarkup       ReplyMarkup `json:"reply_markup,omitempty"`
	SupportsStreaming bool        `json:"supports_streaming"` // pass true if video is suitable for streaming
}

type SendDocumentRequest struct {
	ChatId            int64       `json:"chat_id"`
	Document          string      `json:"document"`
	Caption           string      `json:"caption,omitempty"`
	ParseMode         string      `json:"parse_mode,omitempty"` // for caption
	HasSpoiler        bool        `json:"has_spoiler,omitempty"`
	ReplyMarkup       ReplyMarkup `json:"reply_markup,omitempty"`
	SupportsStreaming bool        `json:"supports_streaming"` // pass true if video is suitable for streaming
}

type SendAnimationRequest struct {
	ChatId      int64       `json:"chat_id"`
	Animation   string      `json:"animation"`
	Caption     string      `json:"caption,omitempty"`
	Duration    int64       `json:"duration,omitempty"` // duration of animation in seconds
	Width       int64       `json:"width"`
	Height      int64       `json:"height"`
	ParseMode   string      `json:"parse_mode,omitempty"` // for caption
	HasSpoiler  bool        `json:"has_spoiler,omitempty"`
	ReplyMarkup ReplyMarkup `json:"reply_markup,omitempty"`
}

type SendVoiceRequest struct {
	ChatId      int64       `json:"chat_id"`
	Voice       string      `json:"voice"`
	Caption     string      `json:"caption,omitempty"`
	Duration    int64       `json:"duration,omitempty"`   // duration of voice message in seconds
	ParseMode   string      `json:"parse_mode,omitempty"` // for caption
	ReplyMarkup ReplyMarkup `json:"reply_markup,omitempty"`
}

type SendVideoNoteRequest struct {
	ChatId      int64       `json:"chat_id"`
	VideoNote   string      `json:"video_note"`
	Duration    int64       `json:"duration,omitempty"`   // duration of sent video in seconds
	ParseMode   string      `json:"parse_mode,omitempty"` // for caption
	ReplyMarkup ReplyMarkup `json:"reply_markup,omitempty"`
}

type SendChatAction struct {
	ChatId int64  `json:"chat_id"`
	Action string `json:"action"` // typing, upload_photo
}
