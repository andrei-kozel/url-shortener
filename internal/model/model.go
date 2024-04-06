package model

import (
	"errors"
	"time"

	. "github.com/samber/mo"
)

var (
	ErrorNotFound                = errors.New("Not found")
	ErrorIdentifierAlreadyExists = errors.New("Identifier already exists")
)

type Shortening struct {
	Identifier  string    `json:"identifier"`
	OriginalUrl string    `json:"original_url"`
	Visits      int64     `json:"visits"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ShorteningInput struct {
	RawUrl     string         `json:"raw_url"`
	Identifier Option[string] `json:"identifier"`
	CreatedBy  string         `json:"created_by"`
}
