package lib

import (
	"fmt"
	"log"
	"net/http"
	"tview_go/config"
)

func SendMessage(msg string) {
	_, err := http.Get(fmt.Sprintf("https://api.telegram.org/"+
		"bot%s/sendMessage?chat_id=%s&text=%s", config.BotStat, config.ChatId, msg))
	if err != nil {
		log.Fatalln(err)
	}
}
