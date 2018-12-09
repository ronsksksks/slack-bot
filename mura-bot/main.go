package main

import (
	"github.com/lob-inc/rssp/server/shared/logger"
	"github.com/nlopes/slack"
	"github.com/riona/slack-bot/mura-bot/schedules"
	"log"
	"os"
)

func main() {
	token := os.Getenv("MURABOT")
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
			case "今週の予定は？":
				schedule, err := schedules.Schedules("week")
				if err != nil {
					logger.Errorf("google-calender error: %v", err)
				}
				rtm.SendMessage(rtm.NewOutgoingMessage(schedule, ev.Channel))
			}
			if ev.Msg.Text == "今日の予定は？" {
				schedule, err := schedules.Schedules("day")
				if err != nil {
					logger.Errorf("google-calender error: %v", err)
				}
				rtm.SendMessage(rtm.NewOutgoingMessage(schedule, ev.Channel))
			}
		case *slack.InvalidAuthEvent:
			log.Print("Invalid credentials")
		}
	}
}
