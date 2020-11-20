package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/mdapathy/url-shortener/domain"
	"github.com/mdapathy/url-shortener/tools"
	"github.com/mdapathy/url-shortener/url/usecase"
	"net/http"
)

type controller struct {
	useCase usecase.UseCase
}

type Controller interface {
	HandleUrlRedirect(rw http.ResponseWriter, r *http.Request)
	HandleUriDelete(rw http.ResponseWriter, r *http.Request)
	HandleUriCreate(rw http.ResponseWriter, r *http.Request)
}

func NewController(u usecase.UseCase) Controller {
	if u == nil {
		panic("Missing use case")
	}
	return &controller{
		u,
	}
}

func (c *controller) HandleUrlRedirect(rw http.ResponseWriter, r *http.Request) {
	key := mux.Vars(r)["key"]
	res, err := c.useCase.GetInitialUrl(key)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusNotFound)
	} else {
		http.Redirect(rw, r, res, http.StatusFound)
	}
}

func (c *controller) HandleUriDelete(rw http.ResponseWriter, r *http.Request) {
	if err := c.useCase.RemoveUrl(mux.Vars(r)["key"]); err != nil {
		http.Error(rw, err.Error(), http.StatusNotFound)
	} else {
		tools.WriteJsonOk(rw, nil)
	}
}

func (c *controller) HandleUriCreate(rw http.ResponseWriter, r *http.Request) {
	var res *domain.Url

	var p domain.UrlDto
	err := json.NewDecoder(r.Body).Decode(&p)

	if err != nil || len(p.Url) == 0 {
		tools.WriteJsonBadRequest(rw, fmt.Sprintf("Error parsing request body"))
		return
	}

	res, err = c.useCase.SaveShortenedUrl(p.Url)
	if res == nil && err != nil {
		tools.WriteJsonBadRequest(rw, err.Error())
	} else {
		tools.WriteJsonCreated(rw, res)
	}
}
