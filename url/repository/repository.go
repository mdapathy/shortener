package repository

import (
	"database/sql"
	"github.com/mdapathy/url-shortener/domain"
)

var (
	insertUrl = "insert into shortened_urls (shortened_url, initial_url) values ($1::varchar , $2::varchar);"
	deleteUrl = "delete from shortened_urls where shortened_url =$1::varchar ;"
	selectUrl = "select shortened_url, initial_url  from  shortened_urls where shortened_url=$1::varchar ;"
)

type repository struct {
	db *sql.DB
}

type Repository interface {
	SaveUrl(*domain.Url) error
	GetUrl(string) (*domain.Url, error)
	RemoveUrl(string) (int64, error)
}

func NewRepository(db *sql.DB) Repository {
	if db == nil {
		panic("Missing database")
	}
	return &repository{
		db: db,
	}
}

func (r *repository) SaveUrl(url *domain.Url) error {
	_, err := r.db.Exec(insertUrl, url.ShortenedUrl, url.InitialUrl)
	return err
}

func (r *repository) GetUrl(shortenedUrl string) (*domain.Url, error) {
	var init, short string
	err := r.db.QueryRow(selectUrl, shortenedUrl).Scan(&short, &init)
	url := domain.Url{
		InitialUrl:   init,
		ShortenedUrl: short,
	}
	return &url, err

}

// Removes url and returns the amount of rows affected
func (r *repository) RemoveUrl(url string) (int64, error) {
	res, err := r.db.Exec(deleteUrl, url)
	if err != nil {
		return 0, err

	}
	return res.RowsAffected()

}
