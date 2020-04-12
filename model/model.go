package model

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"time"
)

const regexExpression = `^([0-9]+):([a-zA-Z]+):([0-9]+|\{.*?\}+)?:?(\{.*\}+|\[\{.*\}\]+)`

const (
	VerbGet       = "GET"
	VerbSet       = "SET"
	VerbAdd       = "ADD"
	VerbRemove    = "REMOVE"
	VerbExtend    = "EXTEND"
	VerbCancel    = "CANCEL"
	VerbSubscribe = "SUBSCRIBE"
	VerbOpen      = "OPEN"
	VerbLoad      = "LOAD"
	VerbReload    = "RELOAD"

	TargetSession      = "session"
	TargetNotification = "notifications"
	TargetHTTP         = "http"
)

// Message ...
type Message struct {
	InvokeID *int
	Verb     string
	Options  interface{}
	Params   interface{}
}

// MessageTarget ...
type MessageTarget struct {
	Target string `json:"target"`
}

// NewMessage ...
func NewMessage(verb string, options interface{}, params interface{}) *Message {
	message := Message{
		Verb:    verb,
		Options: options,
		Params:  params,
	}
	*message.InvokeID = message.GetInvokeID() + 1

	return &message
}

// ToCasper ...
func (msg *Message) ToCasper() []byte {

	return []byte(
		fmt.Sprintf("%d:%s:%v:%v", msg.GetInvokeID(), msg.GetVerb(), msg.GetOptions(), msg.GetParams()),
	)
}

// MessageFromCasper ...
func MessageFromCasper(message []byte) *Message {

	regex := regexp.MustCompile(regexExpression)
	regexRes := regex.FindAllSubmatch(message, -1)

	if regexRes == nil {
		return nil
	}

	msgIndex, _ := strconv.Atoi(
		string(regexRes[0][1]),
	)

	msgParams := regexRes[0][4]

	if regexRes[0][3] != nil && regexRes[0][4] == nil {
		msgParams = regexRes[0][3]
	}

	return &Message{
		InvokeID: &msgIndex,
		Verb:     string(regexRes[0][2]),
		Options:  string(regexRes[0][3]),
		Params:   msgParams,
	}
}

// GetInvokeID ...
func (msg *Message) GetInvokeID() int {
	if msg.InvokeID == nil {
		msg.InvokeID = NewInt(0)
	}

	return *msg.InvokeID
}

// GetVerb ...
func (msg *Message) GetVerb() string {
	return msg.Verb
}

// GetOptions ...
func (msg *Message) GetOptions() interface{} {
	b, _ := json.Marshal(msg.Options)
	return string(b)
}

// GetParams ...
func (msg *Message) GetParams() interface{} {
	b, _ := json.Marshal(msg.Params)
	return string(b)
}

// WaitReceive ...
func WaitReceive(payloadMessage chan *[]byte, messageID int) *Message {
	ticker := time.NewTicker(time.Duration(5) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case msg := <-payloadMessage:
			casperMessage := MessageFromCasper(*msg)

			if casperMessage == nil {
				return nil
			}

			if *casperMessage.InvokeID == messageID {
				return casperMessage
			}
		case <-ticker.C:
			return nil
		}
	}

}
