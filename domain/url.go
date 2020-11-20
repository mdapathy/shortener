package domain

import (
	"github.com/rs/xid"
)

//UrlDto for the post request
type UrlDto struct {
	Url string `json:"url"`
}

//Url entity
type Url struct {
	InitialUrl   string
	ShortenedUrl string
}

//Creates a new entity based on the given url using a unique id generator
func New(initUrl string) *Url {
	return &Url{
		InitialUrl:   initUrl,
		ShortenedUrl: xid.New().String(),
	}
}
