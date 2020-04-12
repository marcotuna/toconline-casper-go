package model

import (
	"encoding/json"
	"io"
)

// Session ...
type Session struct {
	AccessToken string `json:"access_token"`
	InvokeID    int    `json:"invokeId,omitempty"`
	Timer       int    `json:"timer,omitempty"`
}

// ToJSON ...
func (s *Session) ToJSON() string {
	b, _ := json.Marshal(s)
	return string(b)
}

// SessionFromJSON will decode the input and return a Session
func SessionFromJSON(data io.Reader) *Session {
	var session *Session
	json.NewDecoder(data).Decode(&session)
	return session
}
