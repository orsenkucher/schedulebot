package bot

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

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

// GenFor generates keyboard for provided route
func GenFor(route *routeTree) tgbotapi.InlineKeyboardMarkup {
	buttons := make([]tgbotapi.InlineKeyboardButton, len(route.children))
	for i, c := range route.children {
		buttons[i] = tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf(" 📆 %s  ", c.name), "route:"+c.name)
	}
	if route.parent != nil {
		buttons = append([]tgbotapi.InlineKeyboardButton{tgbotapi.NewInlineKeyboardButtonData(" 🔙 Back  ", "back:")}, buttons...)
	}
	return tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(buttons...))
}
