package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/gofrs/uuid"
)

// Device is used by pop to map your .model.Name.Proper.Pluralize.Underscore database table to your go code.
type Device struct {
	ID           uuid.UUID `json:"id" db:"id"`
	Manufacturer string    `json:"manufacturer" db:"manufacturer"`
	Make         string    `json:"make" db:"make"`
	Model        string    `json:"model" db:"model"`
	Storage      string    `json:"storage" db:"storage"`
	Cost         int64     `json:"cost" db:"cost"`
	OS           string    `json:"os" db:"operating_system"`
	ImageURL     string    `json:"image_url" db:"image_url"`
	IsNew        bool      `json:"is_new" db:"is_new"`
	UserID       uuid.UUID `json:"user_id" db:"user_id"`
	User         User      `json:"user" belongs_to:"user"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

type OS []string

// String is not required by pop and may be deleted
func (d Device) String() string {
	jd, _ := json.Marshal(d)
	return string(jd)
}

// Devices is not required by pop and may be deleted
type Devices []Device

// String is not required by pop and may be deleted
func (d Devices) String() string {
	jd, _ := json.Marshal(d)
	return string(jd)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (d *Device) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: d.Manufacturer, Name: "Manufacturer"},
		&validators.StringIsPresent{Field: d.Make, Name: "Make"},
		&validators.StringIsPresent{Field: d.Model, Name: "Model"},
		&validators.StringIsPresent{Field: d.Storage, Name: "Storage"},
		//&validators.FuncValidator{Field: d.Cost, Name: "Cost"},
		&validators.StringIsPresent{Field: d.OS, Name: "OS"},
		&validators.StringIsPresent{Field: d.ImageURL, Name: "ImageURL"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (d *Device) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (d *Device) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
