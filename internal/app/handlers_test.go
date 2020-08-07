package app

import (
	"flag"
	"net/http"
	"os"
	"testing"

	"github.com/rdnply/url-shortener/internal/baseconv"
	"github.com/rdnply/url-shortener/internal/link"
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

	tc := test.TemplateTestCase{
		Name:    "show main page ok",
		Method:  "GET",
		URL:     "/",
		Body:    "",
		Handler: mockApp.loadMainPage,
		Golden:  "main_page",
	}

	test.EndpointReturnsTemplate(t, tc, *update)
}

var (
	counterValue uint = 3
	links             = []*link.Link{
		{1, "example.com", "1", 1, 3},
		{2, "ex.com", "2", 2, 0},
		{3, "examp.org", "3", 3, 4},
	}
)

func TestCreateLink(t *testing.T) {
	mockApp := appForTest()

	mockApp.CounterStorage = &test.MockCounterStorage{Value: counterValue}
	mockApp.LinkStorage = &test.MockLinkStorage{Items: links}

	testCases := []test.EndpointTestCase{
		{"create link ok", "POST", "newLink.com", "application/x-www-form-urlencoded", "url=newLink.com",
			mockApp.createLink, http.StatusMovedPermanently, "", "/stats/4"},
		{"create link get empty string", "GET", "", "", "",
			mockApp.createLink, http.StatusBadRequest, "*get empty string*", ""},
	}

	for _, tc := range testCases {
		test.Endpoint(t, tc)
	}
}

func TestShowStats(t *testing.T) {
	mockApp := appForTest()

	mockApp.LinkStorage = &test.MockLinkStorage{Items: links}

	tc := test.TemplateTestCase{
		Name:    "show stats page ok",
		Method:  "GET",
		URL:     "/stats/1",
		Body:    "",
		Handler: mockApp.showStats,
		Golden:  "stats_page",
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
