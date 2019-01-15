package main

import (
	"github.com/ronsksksks/slack-bot/mura-bot/schedules"
	"log"
	"os"

	"github.com/lob-inc/rssp/server/shared/logger"
	"github.com/nlopes/slack"
)

func main() {
	token := os.Getenv("SLACKBOT")
	api := slack.New(token)

	// connect to Slack RTM with WebSocket
	rtm := api.NewRTM()

	go rtm.ManageConnection()

	// get google calendar events
	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.MessageEvent:
			switch ev.Msg.Text {
			case "今週の予定は？":
				schedule, err := schedules.Schedules("week")
				if err != nil {
					logger.Errorf("google-calender error: %v", err)
				}
				rtm.SendMessage(rtm.NewOutgoingMessage(schedule, ev.Channel))
			case "今日の予定は？":
				schedule, err := schedules.Schedules("day")
				if err != nil {
					logger.Errorf("google-calender error: %v", err)
				}
				rtm.SendMessage(rtm.NewOutgoingMessage(schedule, ev.Channel))
			default:
				rtm.SendMessage(rtm.NewOutgoingMessage("今週の予定は？　や、　今日の予定は？　と聞いてみてね", ev.Channel))
			}
		case *slack.InvalidAuthEvent:
			log.Print("Invalid credentials")
		}
	}
}

