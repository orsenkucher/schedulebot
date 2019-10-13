package job

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func Job(b *tgbotapi.BotAPI) {
	fmt.Println(b.Self.UserName)
}
