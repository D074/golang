package capture

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

var slackClient *SlackClient = &SlackClient{
	TimeOut: 5 * time.Second,
}

func InitSlack(c *SlackClient) {
	slackClient = c
}
func Send(message SlackMessage) error {
	if slackClient.WebHookUrl == "" {
		return errors.New("slack webhook is empty")
	}
	if slackClient.Channel == "" {
		return errors.New("slack channel is empty")
	}

	return sendHttpRequest(message)
}

func sendHttpRequest(slackRequest SlackMessage) error {
	slackBody, _ := json.Marshal(slackRequest)
	req, err := http.NewRequest(http.MethodPost, slackClient.WebHookUrl, bytes.NewBuffer(slackBody))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{Timeout: slackClient.TimeOut}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return err
	}
	if b := buf.String(); b != "ok" {
		return errors.New("non-ok response returned from Slack, got:" + b)
	}
	return nil
}
