package test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/rdnply/url-shortener/internal/project"
	"github.com/rdnply/url-shortener/pkg/pkg/logger"
	"github.com/stretchr/testify/assert"
)

type TemplateTestCase struct {
	Name    string
	Method  string
	URL     string
	Body    string
	Handler http.HandlerFunc
	Golden  string
}

func EndpointReturnsTemplate(t *testing.T, tc TemplateTestCase, update bool) {
	t.Run(tc.Name, func(t *testing.T) {
		req, err := http.NewRequest(tc.Method, tc.URL, bytes.NewBufferString(tc.Body))
		if err != nil {
			t.Fatalf("can't create test request %v", err)
		}

		res := httptest.NewRecorder()

		tc.Handler.ServeHTTP(res, req)

		got := res.Body.String()
		want := goldenValue(t, tc.Golden, got, update)

		assert.Equal(t, got, want, wrongBody(got, want))
	})
}

func goldenValue(t *testing.T, goldenName string, actual string, update bool) string {
	t.Helper()

	os.Chdir(project.Root)
	os.Chdir("test")
	defer os.Chdir(project.Root)

	goldenPath := "testdata/" + goldenName + ".golden"

	f, err := os.OpenFile(goldenPath, os.O_RDWR, 0644)
	assert.Nil(t, err, "can't open golden file: %s", goldenPath)
	defer f.Close()

	if update {
		_, err := f.WriteString(actual)
		assert.Nil(t, err, "error writing to file %s: %s", goldenPath, err)

		return actual
	}

	content, err := ioutil.ReadAll(f)
	assert.Nil(t, err, "error opening file %s: %s", goldenPath, err)

	return string(content)
}

type EndpointTestCase struct {
	Name         string
	Method       string
	URL          string
	Header       string
	Body         string
	Handler      http.HandlerFunc
	WantStatus   int
	WantBody     string
	WantLocation string
}

func Endpoint(t *testing.T, tc EndpointTestCase) {
	t.Run(tc.Name, func(t *testing.T) {
		req, err := http.NewRequest(tc.Method, tc.URL, bytes.NewBufferString(tc.Body))
		if err != nil {
			t.Fatalf("can't create test request %v", err)
		}

		if tc.Header != "" {
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		}

		res := httptest.NewRecorder()

		tc.Handler.ServeHTTP(res, req)
		assert.Equal(t, tc.WantStatus, res.Code, wrongCode(res.Code, tc.WantStatus))

		if tc.WantBody != "" {
			pattern := strings.Trim(tc.WantBody, "*")
			if pattern != tc.WantBody {
				assert.Contains(t, res.Body.String(), pattern, wrongBody(res.Body.String(), pattern))
			} else {
				assert.JSONEq(t, tc.WantBody, res.Body.String(), wrongBody(res.Body.String(), tc.WantBody))
			}
		} else {
			got := res.HeaderMap.Get("Location")
			assert.Equal(t, got, tc.WantLocation, wrongBody(got, tc.WantLocation))
		}

	})
}

func wrongCode(actual, want int) string {
	return fmt.Sprintf("returned wrong status code: got %v, want %v", actual, want)
}

func wrongBody(actual, want string) string {
	return fmt.Sprintf("returned unexpected body: got %v,\n want %v", actual, want)
}

func Logger() logger.Logger {
	config := logger.Configuration{
		EnableConsole:     true,
		ConsoleLevel:      logger.Debug,
		ConsoleJSONFormat: true,
	}

	logger, err := logger.New(config, logger.InstanceZapLogger)
	if err != nil {
		log.Fatal("could not instantiate logger: ", err)
	}

	return logger
}
