package app

import (
	"bytes"
	"fmt"
	"github.com/OB1Company/filehive/repo"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type apiTests []apiTest

type apiTest struct {
	name             string
	path             string
	method           string
	body             []byte
	statusCode       int
	setup            func(db *repo.Database) error
	expectedResponse []byte
}

func errorReturn(err error) []byte {
	return []byte(fmt.Sprintf(`{"error": "%s"}%s`, err.Error(), "\n"))
}

func runAPITests(t *testing.T, tests apiTests) {
	db, err := repo.NewDatabase("", repo.Dialect("test"))
	if err != nil {
		t.Fatal(err)
	}

	server := &FileHiveServer{
		db: db,
	}

	r := server.newV1Router()
	ts := httptest.NewServer(r)
	defer ts.Close()

	var cookies []*http.Cookie
	for _, test := range tests {
		if test.setup != nil {
			if err := test.setup(db); err != nil {
				t.Fatal(err)
			}
		}

		req, err := http.NewRequest(test.method, fmt.Sprintf("%s%s", ts.URL, test.path), bytes.NewReader(test.body))
		if err != nil {
			t.Fatal(err)
		}
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		if res.StatusCode != test.statusCode {
			t.Errorf("%s. Expected status code %d, got %d", test.name, test.statusCode, res.StatusCode)
			continue
		}
		if len(res.Cookies()) > 0 {
			cookies = res.Cookies()
		}

		response, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			log.Fatal(err)
		}

		if test.expectedResponse != nil && !bytes.Equal(response, test.expectedResponse) {
			t.Errorf("%s: Expected response %s, got %s", test.name, string(test.expectedResponse), string(response))
			continue
		}
	}
}
