package capture

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strconv"
	"time"
)

var captureConfig *Config

func InitCapture(c *Config) {
	captureConfig = c
	switch captureConfig.OutputPanic {
	case Slack:
		InitSlack(&c.SlackClient)
	}
}

func CapturePanic(h http.Handler) http.Handler {
	if captureConfig == nil {
		log.Panicln("capture config is not set")
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				switch captureConfig.OutputPanic {
				case Slack:
					notifySlack(err)
				default:
					log.Println(captureConfig.OutputPanic, " has not supported yet")
				}
			}
		}()

		h.ServeHTTP(w, r)
	})
}

func notifySlack(err interface{}) {
	if captureConfig.SlackClient.WebHookUrl == "" {
		log.Println("web hook url for slack cannot be empty")
		return
	}
	if captureConfig.SlackClient.Channel == "" {
		log.Println("channel for slack cannot be empty")
		return
	}
	var attachment Attachment

	text := ""
	if captureConfig.SlackClient.Environment != "" {
		text = "[*Env: " + captureConfig.SlackClient.Environment + "*]\n"
	}
	text += "Panic throws! please check"
	if captureConfig.SlackClient.MentionHere {
		text += "<!here> "
	}
	for _, user := range captureConfig.SlackClient.MentionUser {
		text += "<@" + user + "> "
	}
	text += "\n"
	text += fmt.Sprintf("Error: %s\n", err)
	if captureConfig.ShowTrace {
		buf := make([]byte, 2048)
		n := runtime.Stack(buf, false)
		buf = buf[:n]
		attachment = Attachment{
			Color: "danger",
			Text:  string(buf),
			Ts:    json.Number(strconv.FormatInt(time.Now().Unix(), 10)),
		}
		text += "*Trace:*"
	}
	msg := SlackMessage{
		Username:    captureConfig.SlackClient.UserName,
		Text:        text,
		IconEmoji:   ":zap",
		Attachments: []Attachment{attachment},
		Channel:     captureConfig.SlackClient.Channel,
	}
	err = Send(msg)
	if err != nil {
		log.Println(err)
	}
}
