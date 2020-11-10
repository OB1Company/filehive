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
	expectedResponse func() ([]byte, error)
}

func runAPITests(t *testing.T, tests apiTests) {
	db, err := repo.NewDatabase("", repo.Dialect("test"))
	if err != nil {
		t.Fatal(err)
	}

	server := &FileHiveServer{
		db: db,
	}

	ts := httptest.NewServer(server.newV1Router())
	defer ts.Close()

	for _, test := range tests {
		req, err := http.NewRequest(test.method, fmt.Sprintf("%s%s", ts.URL, test.path), bytes.NewReader(test.body))
		if err != nil {
			t.Fatal(err)
		}
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		if res.StatusCode != test.statusCode {
			t.Errorf("%s. Expected status code %d, got %d", test.name, test.statusCode, res.StatusCode)
			continue
		}
		response, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			log.Fatal(err)
		}
		expected, err := test.expectedResponse()
		if err != nil {
			log.Fatal(err)
		}
		if expected != nil && !bytes.Equal(response, expected) {
			t.Errorf("%s: Expected response %s, got %s", test.name, string(expected), string(response))
			continue
		}
	}
}
