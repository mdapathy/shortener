package middleware

import (
	"database/sql"
	"github.com/mdapathy/url-shortener/tools"
	"github.com/mdapathy/url-shortener/url/controller"
	"github.com/mdapathy/url-shortener/url/repository"
	"github.com/mdapathy/url-shortener/url/usecase"
)

type middleware struct {
	db    *sql.DB
	cache tools.CacheStorage
}

type Middleware interface {
	NewController() controller.Controller
}

func New(db *sql.DB, storage tools.CacheStorage) Middleware {
	if db == nil {
		panic("Missing db")
	}
	return &middleware{
		db:    db,
		cache: storage,
	}
}

func (c *middleware) NewController() controller.Controller {
	return controller.NewController(c.newUseCase())
}

func (c *middleware) newUseCase() usecase.UseCase {
	return usecase.NewUseCase(c.newRepository(), c.cache)
}

func (c *middleware) newRepository() repository.Repository {
	return repository.NewRepository(c.db)
}
