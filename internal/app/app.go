package app

import (
	"html/template"
	"io"
	"log"
	"os"

	"github.com/rdnply/url-shortener/internal/baseconv"
	"github.com/rdnply/url-shortener/internal/counter"
	"github.com/rdnply/url-shortener/internal/link"
	"github.com/rdnply/url-shortener/internal/postgres"
	"github.com/rdnply/url-shortener/pkg/pkg/logger"
)

type App struct {
	Addr           string
	BaseConvertor  *baseconv.BaseConv
	Templates      *templates
	CounterStorage counter.Storage
	LinkStorage    link.Storage
	Logger         logger.Logger
}

func New(addr string) (*App, map[string]io.Closer, error) {
	closers := make(map[string]io.Closer)

	db, err := postgres.New()
	if err != nil {
		return nil, nil, err
	}
	closers["postgres"] = db

	counterStorage, err := postgres.NewCounterStorage(db)
	if err != nil {
		return nil, nil, err
	}
	closers["counter storage"] = counterStorage

	if err := counterStorage.Init(); err != nil {
		return nil, nil, err
	}

	linkStorage, err := postgres.NewLinkStorage(db)
	if err != nil {
		return nil, nil, err
	}
	closers["link storage"] = linkStorage

	baseconv, _ := baseconv.NewBaseConv(62)

	return &App{
		Addr:           addr,
		BaseConvertor:  baseconv,
		Templates:      readTemplates(),
		CounterStorage: counterStorage,
		LinkStorage:    linkStorage,
		Logger:         initLogger(),
	}, nil, nil
}

type templates struct {
	main *template.Template
}

func readTemplates() *templates {
	return &templates{
		main: createTemplate("main", "main.html", "new-form.html", "stats.html"),
	}
}

func createTemplate(name string, tmpls ...string) *template.Template {
	os.Chdir("templates")

	t := template.Must(template.New(name).
		ParseFiles(tmpls...))

	return t
}

func initLogger() logger.Logger {
	config := logger.Configuration{
		EnableConsole:     true,
		ConsoleLevel:      logger.Debug,
		ConsoleJSONFormat: true,
		EnableFile:        true,
		FileLevel:         logger.Info,
		FileJSONFormat:    true,
		FileLocation:      "log.log",
	}

	logger, err := logger.New(config, logger.InstanceZapLogger)
	if err != nil {
		log.Fatal("could not instantiate logger: ", err)
	}

	return logger
}
