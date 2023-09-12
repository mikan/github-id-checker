package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Message defines the chat message.
// Specification: https://developers.google.com/hangouts/chat/reference/message-formats/basic
type Message struct {
	Text string `json:"text"`
}

// Send sends a message to specified webhook url.
func (m *Message) Send(webhook string) error {
	payload, err := json.Marshal(m)
	if err != nil {
		return err
	}
	resp, err := http.Post(webhook, "application/json; charset=UTF-8", bytes.NewReader(payload))
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, body)
	}
	return nil
}
