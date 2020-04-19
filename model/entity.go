package model

import (
	"encoding/json"
	"io"
)

// Entity ...
type Entity struct {
	ID              int    `json:"id"`
	TaxNumber       string `json:"tax_number"`
	Name            string `json:"name"`
	ModuleMask      string `json:"module_mask"`
	RoleMask        string `json:"role_mask"`
	Cluster         string `json:"cluster"`
	Status          string `json:"status"`
	Demo            bool   `json:"demo"`
	CommercialUntil string `json:"commercial_until"`
	AdminUntil      string `json:"admin_util"`
	PastModuleMask  int    `json:"past_module_mask"`
}

// EntityFiscalYears ...
type EntityFiscalYears struct {
	URL         string `json:"url"`
	IssuerURL   string `json:"issuer_url"`
	AccessTTL   int    `json:"access_ttl"`
	UserEmail   string `json:"user_email"`
	EntityID    int    `json:"entity_id"`
	RoleMask    int    `json:"role_mask"`
	ModuleMask  int    `json:"module_mask"`
	AccessToken string `json:"access_token"`
}

// ToJSON convert a Entity to a json string
func (e *Entity) ToJSON() string {
	b, _ := json.Marshal(e)
	return string(b)
}

// EntityListToJSON convert a Entity array to json string
func EntityListToJSON(e []*Entity) string {
	b, _ := json.Marshal(e)
	return string(b)
}

// EntityFromJSON will decode the input and return a Entity
func EntityFromJSON(data io.Reader) *Entity {
	var entity *Entity
	json.NewDecoder(data).Decode(&entity)
	return entity
}

// EntityListFromJSON will decode the input and return a Entity
func EntityListFromJSON(data io.Reader) []*Entity {
	var entities []*Entity
	json.NewDecoder(data).Decode(&entities)
	return entities
}

// EntityFiscalYearsFromJSON will decode the input and return a Entity
func EntityFiscalYearsFromJSON(data io.Reader) *EntityFiscalYears {
	var entityFiscalYears *EntityFiscalYears
	json.NewDecoder(data).Decode(&entityFiscalYears)
	return entityFiscalYears
}
