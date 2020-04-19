package model

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
)

const regexExpression = `^([0-9]+):([a-zA-Z]+):([0-9]+|\{.*?\}+)?:?(\{.*\}+|\[\{.*\}\]+)`

const (
	VerbGet       = "GET"
	VerbPut       = "PUT"
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

// CasperRequest ...
type CasperRequest struct {
	InvokeID *int
	Verb     string
	Options  interface{}
	Params   interface{}
}

// CasperResponse ...
type CasperResponse struct {
	InvokeID *int
	Verb     string
	Options  []byte
	Params   []byte
}

// MessageTarget ...
type MessageTarget struct {
	Target  string      `json:"target"`
	URL     string      `json:"url,omitempty"`
	Headers http.Header `json:"headers,omitempty"`
}

// MessageStatus ...
type MessageStatus struct {
	Success bool `json:"success"`
}

// MessageStatusFromJSON will decode the input and return a MessageStatus
func MessageStatusFromJSON(data io.Reader) *MessageStatus {
	var messageStatus *MessageStatus
	json.NewDecoder(data).Decode(&messageStatus)
	return messageStatus
}

// SetMessage ...
func (msg *CasperRequest) SetMessage(verb string, options interface{}, params interface{}) *CasperRequest {

	message := &CasperRequest{
		Verb:     verb,
		Options:  options,
		Params:   params,
		InvokeID: msg.GetInvokeID(),
	}

	// Increment message id
	*message.InvokeID++

	return message
}

// NewCasperRequest ...
func NewCasperRequest() *CasperRequest {
	return &CasperRequest{InvokeID: NewInt(0)}
}

// ToCasper ...
func (msg *CasperRequest) ToCasper() []byte {

	return []byte(
		fmt.Sprintf("%d:%s:%v:%v", *msg.InvokeID, msg.GetVerb(), msg.GetOptions(), msg.GetParams()),
	)
}

// MessageFromCasper ...
func MessageFromCasper(message []byte) *CasperResponse {

	fmt.Println(string(message))

	regex := regexp.MustCompile(regexExpression)
	regexRes := regex.FindAllSubmatch(message, -1)

	if regexRes == nil {
		return nil
	}

	msgIndex, _ := strconv.Atoi(
		string(regexRes[0][1]),
	)

	msgParams := regexRes[0][4]
	msgOptions := regexRes[0][3]

	if regexRes[0][3] != nil && regexRes[0][4] == nil {
		msgParams = regexRes[0][3]
	}

	return &CasperResponse{
		InvokeID: &msgIndex,
		Verb:     string(regexRes[0][2]),
		Options:  msgOptions,
		Params:   msgParams,
	}
}

// GetInvokeID ...
func (msg *CasperRequest) GetInvokeID() *int {
	if msg.InvokeID == nil {
		msg.InvokeID = NewInt(0)
	}

	return msg.InvokeID
}

// GetVerb ...
func (msg *CasperRequest) GetVerb() string {
	return msg.Verb
}

// GetOptions ...
func (msg *CasperRequest) GetOptions() interface{} {
	b, _ := json.Marshal(msg.Options)
	return string(b)
}

// GetParams ...
func (msg *CasperRequest) GetParams() interface{} {
	b, _ := json.Marshal(msg.Params)
	return string(b)
}
