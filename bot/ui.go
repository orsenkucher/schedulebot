package bot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

var inlineKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(" 📆 1 група  ", "sub:group1"), // 👽 🔴
		tgbotapi.NewInlineKeyboardButtonData(" 📆 2 група  ", "sub:group2"), //👥  🔵 👾 ⏱️
	),
)

var inlineResetKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(" ♻️ 1 група  ", "reset:group1"),
		tgbotapi.NewInlineKeyboardButtonData(" ♻️ 2 група  ", "reset:group2"),
	),
)

var cmdMapping = map[string]string{
	"sub:group1":   "1 група",
	"sub:group2":   "2 група",
	"reset:group1": "1 група",
	"reset:group2": "2 група",
}
