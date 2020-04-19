package app

import (
	"bytes"
	"errors"
	"fmt"
	"toconline-casper-go/model"
	"toconline-casper-go/utils"
)

// SetUserSession ...
func (a *App) SetUserSession(session *model.Authentication) {
	if a.UserSession == nil {
		a.UserSession = session
	}
}

// Authenticate authenticates with username and password
func (a *App) Authenticate(username string, password string) (*model.Authentication, error) {

	reqClient, err := utils.HTTPClientReq(
		"https://app3.toconline.pt/login/sign-in",
		nil,
		[]*utils.HTTPClientHeader{
			{
				Key:   "Authorization",
				Value: fmt.Sprintf("Basic %s", utils.SetBasicAuth(username, password)),
			},
		},
		nil,
	)

	if err != nil {
		return nil, err
	}

	if reqClient.StatusCode < 200 && reqClient.StatusCode > 300 {
		return nil, errors.New("Could not authenticate with the provided credentials")
	}

	userSession := model.AuthenticationFromJSON(bytes.NewReader(reqClient.Body))

	if userSession == nil || len(userSession.AccessToken) == 0 {
		return nil, errors.New("Session is not valid")
	}

	return userSession, nil
}

// SetSession ...
func (a *App) SetSession(userSession *model.Authentication) (*model.Session, error) {
	// Set user session on Casper
	setSession := a.Request.SetMessage(
		model.VerbSet,
		model.MessageTarget{
			Target: model.TargetSession,
		},
		userSession,
	)

	message, err := a.CasperRequest(setSession)

	if err != nil {
		return nil, err
	}

	// Set user session
	a.SetUserSession(userSession)

	return model.SessionFromJSON(bytes.NewReader(message.Params)), nil
}
