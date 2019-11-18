package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/nlopes/slack"
	"github.com/nlopes/slack/slackevents"
)

type EventHandler struct {
	sc *slack.Client
	kc *KibelaClient
}

func NewEventHandler(slackToken, kibelaTeam, kibelaToken string) *EventHandler {
	return &EventHandler{
		sc: slack.New(slackToken),
		kc: NewKibelaClient(kibelaTeam, kibelaToken),
	}
}

func (h *EventHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	body := buf.String()
	eventsAPIEvent, e := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionNoVerifyToken())
	if e != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	switch eventsAPIEvent.Type {
	case slackevents.URLVerification:
		var r *slackevents.ChallengeResponse
		err := json.Unmarshal([]byte(body), &r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(r.Challenge))
	case slackevents.CallbackEvent:
		innerEvent := eventsAPIEvent.InnerEvent
		switch ev := innerEvent.Data.(type) {
		case *slackevents.LinkSharedEvent:
			unfurls := map[string]slack.Attachment{}
			for _, url := range ev.Links {
				ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
				defer cancel()
				res, err := h.kc.NoteFromPath(ctx, url.URL)
				if err != nil {
					log.Printf("[ERROR] can not retrieve the note (%s): %s", url, err)
					return
				}
				t, err := time.Parse(KibelaTimeFormat, res.Note.PublishedAt)
				if err != nil {
					log.Printf("[ERROR] can not parse publishedAt (%s): %s", res.Note.PublishedAt, err)
					return
				}
				unfurls[url.URL] = slack.Attachment{
					AuthorIcon: res.Note.Author.AvatarImage.URL,
					AuthorLink: res.Note.Author.URL,
					AuthorName: res.Note.Author.Account,
					Title:      res.Note.Title,
					TitleLink:  res.Note.URL,
					Text:       res.Note.Summary,
					Footer:     "Kibela",
					FooterIcon: "https://cdn.kibe.la/assets/shortcut_icon-99b5d6891a0a53624ab74ef26a28079e37c4f953af6ea62396f060d3916df061.png",
					Ts:         json.Number(strconv.FormatInt(t.Unix(), 10)),
				}
			}
			options := []slack.MsgOption{}
			_, _, _, err := h.sc.UnfurlMessage(ev.Channel, ev.MessageTimeStamp.String(), unfurls, options...)
			if err != nil {
				log.Printf("[ERROR] can not send unfurl message: %s", err)
				return
			}
			log.Printf("[INFO] send unfurl message for %s", ev.Links)
		}
	}
}
