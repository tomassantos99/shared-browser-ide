package ws

import (
	"encoding/json"
	"errors"
	"slices"
)

type Message struct {
	Type                string `json:"type"`
	ProgrammingLanguage string `json:"programmingLanguage"`
	Content             string `json:"content"`
}

const (
	ClientCodeUpdate  = "ClientCodeUpdate"
	SessionCodeUpdate = "SessionCodeUpdate"
	ClientsUpdate     = "ClientsUpdate"
	Unknown           = "Unknown"
)

var validMessageTypes []string = []string{
	ClientCodeUpdate,
	SessionCodeUpdate,
	ClientsUpdate,
	Unknown,
}

func DefaultMessage() Message {
	return Message{
		Type:    Unknown,
		Content: "",
	}
}

func CreateMessage(messageType string, programmingLanguage string, content string) (Message, error) {
	var message = Message{
		messageType,
		programmingLanguage,
		content,
	}

	var valError = message.Validate()
	if valError != nil {
		return DefaultMessage(), valError
	}

	return message, nil
}

func UnmarshalMessage(bytes []byte) (Message, error) {
	var convertedMessage Message

	var err = json.Unmarshal(bytes, &convertedMessage)
	if err != nil {
		return DefaultMessage(), err
	}

	var valErr = convertedMessage.Validate()
	if valErr != nil {
		return DefaultMessage(), err
	}

	return convertedMessage, nil
}

func (message *Message) Validate() error {
	if !slices.Contains(validMessageTypes, message.Type) {
		return errors.New("invalid message type")
	}
	return nil
}
