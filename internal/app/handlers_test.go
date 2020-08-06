package app

import (
	"flag"
	"os"
	"testing"

	"github.com/rdnply/url-shortener/internal/baseconv"
	"github.com/rdnply/url-shortener/test"
)

// for update or generate golden files run: go test ./... -update
// https://ieftimov.com/post/testing-in-go-golden-files/
var update = flag.Bool("update", false, "update the golden files of this test")

func TestMain(m *testing.M) {
	flag.Parse()
	os.Exit(m.Run())
}

func TestLoadMainPage(t *testing.T) {
	mockApp := appForTest()

	payload := struct {
		NewForm bool
	}{
		NewForm: true,
	}

	tc := test.TemplateTestCase{
		Name:    "show main page ok",
		Method:  "GET",
		URL:     "/",
		Body:    "",
		Handler: mockApp.loadMainPage,
		Payload: payload,
		Golden:  "main_page",
	}

	test.EndpointReturnsTemplate(t, tc, *update)
}

func appForTest() *App {
	baseconv, _ := baseconv.NewBaseConv(62)

	return &App{
		BaseConvertor: baseconv,
		Templates:     readTemplates(),
		Logger:        test.Logger(),
	}
}
