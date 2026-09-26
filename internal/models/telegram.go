package models

type EmojiReactionPayload struct {
	Type  string `json:"type"`  // всегда "emoji"
	Emoji string `json:"emoji"` // emoji-символ
}

type SendOption func(config *SendOptions)

type SendOptions struct {
	ReplyToMessageID int
	SourceName       string
	SourceURL        string
	Caption          string
}

func WithSource(name, url string) SendOption {
	return func(cfg *SendOptions) {
		cfg.SourceName = name
		cfg.SourceURL = url
	}
}

func WithReplyTo(messageID int) SendOption {
	return func(cfg *SendOptions) {
		cfg.ReplyToMessageID = messageID
	}
}

func WithCaption(caption string) SendOption {
	return func(cfg *SendOptions) {
		cfg.Caption = caption
	}
}
