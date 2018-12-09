package main

import (
	"github.com/nlopes/slack"
	"log"
	"os"
	"strings"
)

func main() {
	// API Clientを作成する
	token := os.Getenv("SLACKBOT")
	api := slack.New(token)

	// WebSocketでSlack RTM APIに接続する
	rtm := api.NewRTM()
	// goroutineで並列化する
	go rtm.ManageConnection()

	// イベントを取得する
	for msg := range rtm.IncomingEvents {
		// 型swtichで型を比較する
		switch ev := msg.Data.(type) {
		case *slack.MessageEvent:
			switch ev.Msg.Text {
			case "好きって100回言って":
				var messages []string
				for i := 0; i < 100; i ++ {
					messages = append(messages, "好き")
				}
				rtm.SendMessage(rtm.NewOutgoingMessage(strings.Join(messages, "\n"), ev.Channel))
			case `好き * 100`, `好き * １00`, `すき * 100`, `すき * １００`:
				rtm.SendMessage(rtm.NewOutgoingMessage("怒るよ", ev.Channel))
			case "ジャンプ買ってきたよ", "じゃんぷかってきたよ", "ジャンプかってきたよ", "ジャンプ":
				rtm.SendMessage(rtm.NewOutgoingMessage("やったあ！", ev.Channel))
			case "wifi", "WIFI", "ワイファイ", "わいふぁい":
				rtm.SendMessage(rtm.NewOutgoingMessage("3701012347512", ev.Channel))
			}

		case *slack.InvalidAuthEvent:
			log.Print("Invalid credentials")
		}
	}
}
