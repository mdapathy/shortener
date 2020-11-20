package usecase

import (
	"fmt"
	"github.com/mdapathy/url-shortener/domain"
	"github.com/mdapathy/url-shortener/tools"
	repository "github.com/mdapathy/url-shortener/url/repository"
	"log"
)

type useCase struct {
	repository repository.Repository
	cache      tools.CacheStorage
}

type UseCase interface {
	SaveShortenedUrl(string) (*domain.Url, error)
	GetInitialUrl(string) (string, error)
	RemoveUrl(string) error
}

func NewUseCase(repository repository.Repository, cache tools.CacheStorage) UseCase {
	if repository == nil {
		panic("Missing repository")
	}
	return &useCase{
		repository: repository,
		cache:      cache,
	}
}

// Communicates with repository to save url and returns the url entity
func (u *useCase) SaveShortenedUrl(url string) (*domain.Url, error) {
	urlEntity := domain.New(url)
	err := u.repository.SaveUrl(urlEntity)
	if err != nil {
		return nil, err
	}

	if err := u.cache.SetValue(urlEntity.ShortenedUrl, urlEntity.InitialUrl); err != nil {
		log.Printf("couldn't store `%s` in cache\n", urlEntity.ShortenedUrl)
	}

	log.Printf("stored `%s` in cache\n", urlEntity.ShortenedUrl)

	return urlEntity, nil

}

func (u *useCase) GetInitialUrl(url string) (string, error) {
	str, err := u.cache.GetValue(url)
	if err == nil {
		log.Printf("retrieved value `%s` from cache\n", str)
		return str, nil
	}

	res, _ := u.repository.GetUrl(url)

	if res == nil || len(res.InitialUrl) == 0 {
		return "", fmt.Errorf("Shortened url for `%s` doesn't exist\n", url)
	}

	if err := u.cache.SetValue(res.ShortenedUrl, res.InitialUrl); err != nil {
		log.Printf("couldn't store `%s` in cache\n", res.ShortenedUrl)
	}

	return res.InitialUrl, nil
}
func (u *useCase) RemoveUrl(url string) error {
	removed, err := u.repository.RemoveUrl(url)
	if err != nil {
		return err
	}

	if removed < 1 {
		log.Printf(" value `%s` does not exist \n", url)
		return fmt.Errorf("Url %s doesn't exist\n", url)
	}

	u.cache.DeleteValue(url)
	log.Printf("removed value `%s` from cache\n", url)

	return nil
}
