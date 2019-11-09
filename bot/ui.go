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
		buttons[i] = tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf(" 📆 %s  ", ch.Name), "route:"+ch.MakePath())
	}
	if route.Parent != nil {
		buttons = append([]tgbotapi.InlineKeyboardButton{tgbotapi.NewInlineKeyboardButtonData(" 🔙 Back  ", "route:"+route.Parent.MakePath())}, buttons...)
	}
	return tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(buttons...)), true
}
