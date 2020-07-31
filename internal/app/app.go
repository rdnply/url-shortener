package app

import (
	"html/template"
	"log"
	"os"

	"github.com/rdnply/url-shortener/internal/counter"
	"github.com/rdnply/url-shortener/internal/link"
	"github.com/rdnply/url-shortener/pkg/pkg/logger"
)

type App struct {
	Addr           string
	Templates      *templates
	LinkStorage    link.Storage
	CounterStorage counter.Storage
	Logger         logger.Logger
}

func New(addr string) (*App, error) {
	return &App{
		Addr:      addr,
		Templates: readTemplates(),
		Logger:    initLogger(),
	}, nil
}

type templates struct {
	main *template.Template
}

func readTemplates() *templates {
	return &templates{
		main: createTemplate("main", "main.html", "new-form.html"),
	}
}

func createTemplate(name string, tmpls ...string) *template.Template {
	os.Chdir("../..")
	os.Chdir("static/templates")

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
