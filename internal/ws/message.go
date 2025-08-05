package ws

import (
	"encoding/json"
	"errors"
	"github.com/sirupsen/logrus"
)

type Message struct {
	Type                MessageType `json:"type"`
	ProgrammingLanguage string      `json:"programmingLanguage"`
	Content             string      `json:"content"`
}

type MessageType struct {
	label string
}

const (
	ClientCodeUpdate  = "ClientCodeUpdate"
	SessionCodeUpdate = "SessionCodeUpdate"
	ClientsUpdate     = "ClientsUpdate"
	Unknown           = "Unknown"
)

var validMessageTypes map[string]MessageType = map[string]MessageType{
	ClientCodeUpdate:  {ClientCodeUpdate},
	SessionCodeUpdate: {SessionCodeUpdate},
	ClientsUpdate:     {ClientsUpdate},
	Unknown:           {Unknown},
}

func (messageType MessageType) String() string {
	return messageType.label
}

func messageTypeFromString(messageType string) (MessageType, error) {
	var mType, ok = validMessageTypes[messageType]
	if !ok {
		return validMessageTypes[Unknown], errors.New("Unknown Message Type: " + messageType)
	}

	return mType, nil
}

func DefaultMessage() Message {
	return Message{
		Type:    MessageType{Unknown},
		Content: "",
	}
}

func CreateMessage(messageType string, programmingLanguage string, content string) (Message, error) {
	var mType, err = messageTypeFromString(messageType)
	if err != nil {
		return DefaultMessage(), err
	}

	return Message{
		mType,
		programmingLanguage,
		content,
	}, nil
}

func (messageType MessageType) MarshalJSON() ([]byte, error) {
	return json.Marshal(messageType.label)
}

func (messageType *MessageType) UnmarshalJSON(data []byte) error {
	var typeName string

	var jsonErr = json.Unmarshal(data, &typeName)
	if jsonErr != nil {
		logrus.Error(jsonErr)
	}

	var mType, err = messageTypeFromString(typeName)
	if err != nil {
		logrus.Error(err)
		*messageType = validMessageTypes[Unknown]
		return err
	}

	*messageType = mType
	return nil
}
