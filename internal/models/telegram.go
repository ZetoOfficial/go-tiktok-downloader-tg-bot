package models

type EmojiReactionPayload struct {
	Type  string `json:"type"`  // всегда "emoji"
	Emoji string `json:"emoji"` // emoji-символ
}

type SendOption func(config *SendOptions)

type SendOptions struct {
	ReplyToMessageID int
}

func WithReplyTo(messageID int) SendOption {
	return func(cfg *SendOptions) {
		cfg.ReplyToMessageID = messageID
	}
}
