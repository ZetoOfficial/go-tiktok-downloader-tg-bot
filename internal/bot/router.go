package bot

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func handleCommand(botAPI *tgbotapi.BotAPI, message *tgbotapi.Message) {
	// Устанавливаем доступные команды бота
	commands := []tgbotapi.BotCommand{
		{Command: "start", Description: "Start the bot"},
		{Command: "help", Description: "How to use the bot"},
	}

	_, err := botAPI.Request(tgbotapi.NewSetMyCommands(commands...))
	if err != nil {
		log.Fatalf("Failed to set bot commands: %v", err)
	}

	chatID := message.Chat.ID
	switch message.Command() {
	case "start":
		startText := "👋 *Hi! I’m your video downloader.*\n" +
			"I can help you download videos from *TikTok* and *YouTube Shorts* — fast, free, and without watermarks.\n\n" +
			"🎬 *How to use me:*\n" +
			"1. Send me a TikTok or YouTube Shorts link\n" +
			"2. I’ll fetch and send it back to you!\n\n" +
			"✅ Just send a link to try it out!\n\n" +
			"👥 *Want to use me in a group?*\n" +
			"Add me to the group and give me *admin rights* with the following permissions:\n" +
			"- Read messages\n" +
			"- Send messages\n" +
			"- (Optional) Delete messages — useful for cleanup\n\n" +
			"Without admin rights, I won’t be able to respond to messages in the group."

		msg := tgbotapi.NewMessage(chatID, startText)
		msg.ParseMode = "Markdown"
		_, _ = botAPI.Send(msg)

	case "help":
		helpText := "ℹ️ *Bot Commands:*\n\n" +
			"- Send a TikTok or YouTube Shorts link — I’ll download it for you.\n" +
			"- Videos are sent *without watermarks* when possible.\n\n" +
			"⚠️ Note: Telegram limits video uploads to ~50MB for bots.\n\n" +
			"👥 *Using the bot in a group?*\n" +
			"Make sure I have *admin rights*:\n" +
			"- Read messages\n" +
			"- Send messages\n" +
			"- (Optional) Delete messages\n\n" +
			"I won’t work properly in group chats without these permissions."

		msg := tgbotapi.NewMessage(chatID, helpText)
		msg.ParseMode = "Markdown"
		_, _ = botAPI.Send(msg)

	default:
		msg := tgbotapi.NewMessage(chatID, "❓ Unknown command. Type /help for available options.")
		_, _ = botAPI.Send(msg)
	}
}
