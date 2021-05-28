package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/atreya2011/slack"
	"github.com/atreya2011/slack/slackevents"
)

var api = slack.New("slack-token")

func main() {
	http.HandleFunc("/events", func(w http.ResponseWriter, r *http.Request) {
		// read body
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		// parse event
		eventsAPIEvent, err := slackevents.ParseEvent(
			json.RawMessage(body),
			slackevents.OptionVerifyToken(
				&slackevents.TokenComparator{
					VerificationToken: "gIgiCOBigVnicjmeOm8ONDMQ",
				}),
		)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		// handle url verification event
		if eventsAPIEvent.Type == slackevents.URLVerification {
			var r *slackevents.ChallengeResponse
			err := json.Unmarshal(body, &r)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
			w.Header().Set("Content-Type", "text")
			if _, err := w.Write([]byte(r.Challenge)); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}
		// handle callback event
		if eventsAPIEvent.Type == slackevents.CallbackEvent {
			innerEvent := eventsAPIEvent.InnerEvent
			switch ev := innerEvent.Data.(type) {
			case *slackevents.ReactionAddedEvent:
				_, _, channelID, err := api.OpenIMChannel("UKBTQHYG2")
				if err != nil {
					log.Println(err)
				}
				if ev.ItemUser == "UKBTQHYG2" {
					var channelName, userName, reactedMessage string
					channel, err := api.GetConversationInfo(ev.Item.Channel, false)
					if err == nil {
						channelName = channel.Name
					}
					user, err := api.GetUserInfo(ev.User)
					if err == nil {
						userName = user.Name
					}
					history, _ := api.GetChannelHistory(
						ev.Item.Channel,
						slack.HistoryParameters{
							Latest:    ev.Item.Timestamp,
							Oldest:    ev.Item.Timestamp,
							Inclusive: true,
						})
					if len(history.Messages) > 0 {
						reactedMessage = history.Messages[0].Text
					}
					msg := fmt.Sprintf("%s reacted with :%s: to your message: `%s`, that you posted in `%s`", userName, ev.Reaction, reactedMessage, channelName)
					_, err = api.PostEphemeral(channelID, "UKBTQHYG2", slack.MsgOptionText(msg, true))
					if err != nil {
						log.Println(err)
					}
				}
			}
		}
	})
	log.Println("[INFO] Server listening")
	log.Fatal(http.ListenAndServe(":5002", nil))
}
