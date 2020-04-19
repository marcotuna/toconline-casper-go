package model

import (
	"encoding/json"
	"io"
)

// Authentication ...
type Authentication struct {
	AccessToken string `json:"access_token"`
	InvokeID    int    `json:"invokeId,omitempty"`
	Timer       int    `json:"timer,omitempty"`
}

// Session ...
type Session struct {
	App       SessionApp `json:"app"`
	RoleMask  int        `json:"role_mask"`
	Success   bool       `json:"success"`
	UserEmail string     `json:"user_email"`
	UserID    int        `json:"user_id"`
	UserName  string     `json:"user_name"`
}

// SessionApp ...
type SessionApp struct {
	AdminUntil              string        `json:"admin_until"`
	C4Connected             bool          `json:"c4_connected"`
	CertifiedSoftwareNotice string        `json:"certified_software_notice"`
	CommercialUntil         string        `json:"commercial_until"`
	Config                  SessionConfig `json:"config"`
	Entity                  SessionEntity `json:"entity"`
	IsDemo                  bool          `json:"is_demo"`
}

// SessionConfig ...
type SessionConfig struct {
	CdbAPI                string `json:"cdb_api"`
	JwtJobURL             string `json:"jwt_job_url"`
	PublicAssetsURL       string `json:"public_assets_url"`
	SuperUser             bool   `json:"super_user"`
	UploadURL             string `json:"upload_url"`
	UserBehavesAsEmployee bool   `json:"user_behaves_as_employee"`
}

// SessionEntity ...
type SessionEntity struct {
	AccountantTaxRegistrationNumber string `json:"accountant_tax_registration_number"`
	ActivateAccountingExport        bool   `json:"activate_accounting_export"`
	SocialSecurityNumber            string `json:"social_security_number"`
	TaxRegistrationNumber           string `json:"tax_registration_number"`
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

// ToJSON ...
func (s *Authentication) ToJSON() string {
	b, _ := json.Marshal(s)
	return string(b)
}

// AuthenticationFromJSON will decode the input and return a Authentication
func AuthenticationFromJSON(data io.Reader) *Authentication {
	var authentication *Authentication
	json.NewDecoder(data).Decode(&authentication)
	return authentication
}
