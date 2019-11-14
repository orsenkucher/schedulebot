package bot

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/orsenkucher/schedulebot/route"
)

// GenFor generates keyboard for provided route
func GenFor(route *route.Tree) (tgbotapi.InlineKeyboardMarkup, bool) {
	if route.Children == nil {
		return tgbotapi.InlineKeyboardMarkup{}, false
	}
	buttons := make([]tgbotapi.InlineKeyboardButton, len(route.Children))
	for i, ch := range route.Children {
		icon := "📂"
		if ch.Children == nil {
			icon = "📆"
		}
		buttons[i] = tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf(" %s %s  ", icon, ch.Name), "route:"+ch.CalcHash64())
	}
	if route.Parent != nil {
		buttons = append([]tgbotapi.InlineKeyboardButton{tgbotapi.NewInlineKeyboardButtonData(" 🔙 Back  ", "route:"+route.Parent.CalcHash64())}, buttons...)
	}
	return tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(buttons...)), true
}

// GenForReset generates keyboard for provided reset tree
func GenForReset(route *route.Tree) (tgbotapi.InlineKeyboardMarkup, bool) {
	if route.Children == nil {
		return tgbotapi.InlineKeyboardMarkup{}, false
	}
	buttons := make([]tgbotapi.InlineKeyboardButton, len(route.Children))
	for i, ch := range route.Drop().Children {
		icon := "📂"
		if ch.Children == nil {
			icon = "♻️"
		}
		buttons[i] = tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf(" %s %s  ", icon, ch.Name), "reset:"+ch.CalcHash64())
	}
	if route.Parent != nil {
		buttons = append([]tgbotapi.InlineKeyboardButton{tgbotapi.NewInlineKeyboardButtonData(" 🔙 Back  ", "reset:"+route.Parent.CalcHash64())}, buttons...)
	}
	return tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(buttons...)), true
}
