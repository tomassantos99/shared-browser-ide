package ws

import (
	"encoding/json"
	"errors"
	"github.com/sirupsen/logrus"
)

type Message struct {
	Type MessageType `json:"type"`
	Data string `json:"data"`
}

type MessageType struct {
	label string
}

const (
	clientCodeUpdate= "ClientCodeUpdate"
	sessionCodeUpdate = "SessionCodeUpdate"
	clientsUpdate = "ClientsUpdate"
	unknown = "Unknown"
) 


var validMessageTypes map[string]MessageType = map[string]MessageType {
	clientCodeUpdate: {clientCodeUpdate},
	sessionCodeUpdate: {sessionCodeUpdate},
	clientsUpdate: {clientsUpdate},
	unknown: {unknown},
}

func (messageType MessageType) String() string {
	return messageType.label
}

func messageTypeFromString(messageType string) (MessageType, error){
	var mType, ok = validMessageTypes[messageType]
	if !ok {
		return validMessageTypes[unknown], errors.New("Unknown Message Type: " + messageType)
	}

	return mType, nil
}

func DefaultMessage() Message {
	return Message{
		Type: MessageType{unknown},
		Data: "",
	}
}

func (messageType MessageType) MarshalJSON() ([]byte, error) {
	return json.Marshal(messageType.label)
}

func (messageType *MessageType) UnmarshalJSON(data []byte) error {
	var typeName string
	
	var jsonErr = json.Unmarshal(data, &typeName); if jsonErr != nil {
		logrus.Error(jsonErr)
	}
	
	var mType, err = messageTypeFromString(typeName)
	if err != nil {
		logrus.Error(err)
		*messageType = validMessageTypes[unknown]
		return err
	}

	*messageType = mType
	return nil
}
