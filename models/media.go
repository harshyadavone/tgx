package models

// Photo animation media etc.

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
	MimeType     string `json:"mime_type"`
	FileSize     int64  `json:"file_size"`
}

type Document struct {
	FileId       string `json:"file_id"`
	FileUniqueId string `json:"file_unique_id"`
	FileName     string `json:"file_name"`
	MimeType     string `json:"mime_type"`
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
	MimeType     string `json:"mime_type"`
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
	MimeType     string `json:"mime_type"`
	FileSize     int64  `json:"file_size"`
}
