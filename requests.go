package tgx

import "github.com/harshyadavone/tgx/models"

type SendMessageRequest struct {
	ChatId      int64        `json:"chat_id"`              // Required
	Text        string       `json:"text"`                 // Required
	ParseMode   string       `json:"parse_mode,omitempty"` // MarkdownV2 || HTML
	ReplyMarkup *ReplyMarkup `json:"reply_markup,omitempty"`
	ReplyParams *ReplyParam  `json:"reply_paramaters,omitempty"`
}

type ReplyParam struct {
	MessageId int64 `json:"message_id"` // Required
	ChatId    int64 `json:"chat_id,omitempty"`
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

type SendChatActionRequest struct {
	ChatId int64  `json:"chat_id"`
	Action string `json:"action"` // typing, upload_photo
}

type EditMessageTextRequest struct {
	ChatId      int64                       `json:"chat_id"`
	MessageId   int64                       `json:"message_id"`
	Text        string                      `json:"text"`
	ReplyMarkup models.InlineKeyboardMarkup `json:"reply_markup"`
}

type DeleteMessageRequest struct {
	ChatId    int64 `json:"chat_id"`
	MessageId int64 `json:"message_id"`
}
