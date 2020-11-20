package tests

import (
	"bytes"
	"fmt"
	"github.com/mdapathy/url-shortener/database"
	srvpkg "github.com/mdapathy/url-shortener/server"
	"github.com/mdapathy/url-shortener/tools"
	"github.com/mdapathy/url-shortener/url/middleware"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

type testCase struct {
	input      string
	httpStatus int
}

var (
	handler http.Handler

	URL = "http://localhost:8080"

	testCasesPost = []testCase{
		{
			input:      `{"url" : "github.com"}`, // this value should be removed from the db after the test run or, alternatively, changed
			httpStatus: 201,
		},
		{
			input:      `{"url" : "https://www.google.com/"}`,
			httpStatus: 400, // cannot create another link for a pre saved link
		},
		{
			input:      ``,
			httpStatus: 400,//input should be given to create a link
		},
	}

	testCasesGet = []testCase{
		{
			input:      `burmkvfjcu8e4vgummj0`, //shortened url corresponds to a pre saved url
			httpStatus: 200,
		},
		{
			input:      `burg123456`,
			httpStatus: 404,
		},
	}

	testCasesDelete = []testCase{
		{
			input:      `burggtnjcu81j3f50hgg`, // should be changed for the next run
			httpStatus: 200,
		},
		{
			input:      `newvalue`,
			httpStatus: 404,
		},
	}
)

//Setting up the server the way it is done in main.go
func TestMain(m *testing.M) {
	db := database.NewDBConfig("config_test.json")
	defer db.Close()

	mid := middleware.New(db, tools.NewCacheStorage())

	server := &srvpkg.ApiServer{
		Port:       8080,
		Controller: mid.NewController(),
	}

	handler = server.Create()
	srv := httptest.NewServer(handler)
	code := m.Run()
	defer srv.Close()
	os.Exit(code)
}

func TestServerPost(t *testing.T) {
	for _, tCase := range testCasesPost {
		res, err := http.Post(fmt.Sprintf("%s/shorten", URL), "application/json", bytes.NewBuffer([]byte(tCase.input)))
		assert.Nil(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, tCase.httpStatus, res.StatusCode)

	}
}

func TestServerGet(t *testing.T) {
	for _, tCase := range testCasesGet {
		res, err := http.Get(fmt.Sprintf("%s/url/%s", URL, tCase.input))
		assert.Nil(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, tCase.httpStatus, res.StatusCode)

	}

}

func TestServerDelete(t *testing.T) {
	client := &http.Client{}

	for _, tCase := range testCasesDelete {
		req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/url/%s", URL, tCase.input), nil)
		if err != nil {
			t.Errorf("error creating request:%s", err)
		}

		res, err := client.Do(req)

		assert.Nil(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, tCase.httpStatus, res.StatusCode)

	}

}
