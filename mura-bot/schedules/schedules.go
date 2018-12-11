package schedules

import (
	"fmt"
	"github.com/ronsksksks/slack-bot/mura-bot/shared/google_calender"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
)

type CalenderService struct {
	Service *calendar.Service
}
// service struct wrapper
func newService(client *google_calender.GoogleCalender) (*CalenderService, error) {
	srv, err := calendar.New(client.Client)
	return &CalenderService{Service: srv}, err
}

// execute time parse for week or day and
// get the data for the specified time
func Schedules(duration string) (string, error) {
	client, err := Client()
	if err != nil {
		return "", err
	}

	srv, err := newService(client)
	if err != nil {
		return "", err
	}

	now := time.Now()
	var end time.Time
	switch duration {
	case "day":
		end = now.AddDate(0, 0, 1)
	case "week":
		end = now.AddDate(0, 0, 7)
	}
	if end.IsZero() {
		end = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	}
	events, err := srv.Event(now, end)
	if err != nil {
		return "", err
	}

	var messages []string
	if len(events.Items) == 0 {
		fmt.Println("No upcoming events found.")
	} else {
		for _, item := range events.Items {
			date := item.Start.DateTime
			if date == "" {
				date = item.Start.Date
			}
			messages = append(messages, fmt.Sprintf("%v: %v\n", date, item.Summary))
		}
	}
	return strings.Join(messages, ""), nil
}


// google calendar client
func Client() (*google_calender.GoogleCalender, error) {
	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
		return nil, err
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, calendar.CalendarReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client, err := google_calender.Client(config)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// get the data for the specified time
// from google calendar
func (srv *CalenderService) Event(start, end time.Time) (*calendar.Events, error) {
	events, err := srv.Service.Events.List("5t9khdis149i45g9a43dvg7vn8@group.calendar.google.com").
		ShowDeleted(false).SingleEvents(true).TimeMin(start.Format(time.RFC3339)).
		MaxResults(10).TimeMax(end.Format(time.RFC3339)).OrderBy("startTime").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve next ten of the user's events: %v", err)
		return nil, err
	}
	return events, nil
}
