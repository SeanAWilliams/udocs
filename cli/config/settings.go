package config

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"strings"

	"github.com/ultimatesoftware/udocs/cli/udocs"
)

type Settings struct {
	APIVersion        string
	EntryPoint        string
	BindAddr          string
	Port              string
	RootRoute         string
	Organization      string
	SearchPlaceholder string
	Email             string
	Routes            []string
	Seed              []string
	MongoURL          string
	LogoURL           string
	QuipAccessToken   string
}

func LoadSettings() Settings {
	return EnvVars(Conf())
}

func DefaultSettings() Settings {
	return Settings{
		APIVersion:        "v1",
		EntryPoint:        "http://localhost",
		BindAddr:          "0.0.0.0",
		Port:              "9554",
		SearchPlaceholder: "Search",
		Routes:            []string{},
		Seed:              []string{},
	}
}

func (s Settings) String() string {
	buf := new(bytes.Buffer)
	buf.WriteString("UDOCS_API_VERSION=" + s.APIVersion)
	buf.WriteString("\nUDOCS_ENTRY_POINT=" + s.EntryPoint)
	buf.WriteString("\nUDOCS_BIND_ADDR=" + s.BindAddr)
	buf.WriteString("\nUDOCS_PORT=" + s.Port)
	buf.WriteString("\nUDOCS_ROOT_ROUTE=" + s.RootRoute)
	buf.WriteString("\nUDOCS_ROUTES=" + sliceToString(s.Routes))
	buf.WriteString("\nUDOCS_SEED=" + sliceToString(s.Seed))
	buf.WriteString("\nUDOCS_MONGO_URL=" + s.MongoURL)
	buf.WriteString("\nUDOCS_LOGO_URL=" + s.LogoURL)
	buf.WriteString("\nUDOCS_QUIP_ACCESS_TOKEN=" + s.QuipAccessToken)
	return buf.String()
}

func EnvVars(settings Settings) Settings {
	return Settings{
		APIVersion:        loadEnvVar("UDOCS_API_VERSION", settings.APIVersion),
		EntryPoint:        loadEnvVar("UDOCS_ENTRY_POINT", settings.EntryPoint),
		BindAddr:          loadEnvVar("UDOCS_BIND_ADDR", settings.BindAddr),
		Port:              loadEnvVar("UDOCS_PORT", settings.Port),
		RootRoute:         loadEnvVar("UDOCS_ROOT_ROUTE", settings.RootRoute),
		Organization:      loadEnvVar("UDOCS_ORGANIZATION", settings.Organization),
		Email:             loadEnvVar("UDOCS_EMAIL", settings.Email),
		SearchPlaceholder: loadEnvVar("UDOCS_SEARCH_PLACEHOLDER", settings.SearchPlaceholder),
		Routes:            strings.Split(loadEnvVar("UDOCS_ROUTES", sliceToString(settings.Routes)), ","),
		Seed:              strings.Split(loadEnvVar("UDOCS_SEED", sliceToString(settings.Seed)), ","),
		MongoURL:          loadEnvVar("UDOCS_MONGO_URL", settings.MongoURL),
		LogoURL:           loadEnvVar("UDOCS_LOGO_URL", settings.LogoURL),
		QuipAccessToken:   loadEnvVar("UDOCS_QUIP_ACCESS_TOKEN", settings.QuipAccessToken),
	}
}

func loadEnvVar(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		value = defaultValue
	}
	return value
}

func Conf() Settings {
	f, err := os.Open(udocs.ConfPath())
	if err != nil {
		return EnvVars(DefaultSettings())
	}
	defer f.Close()
	return loadFromMap(toMap(f))
}

func loadFromMap(m map[string]string) Settings {
	return Settings{
		APIVersion:        m["UDOCS_API_VERSION"],
		EntryPoint:        m["UDOCS_ENTRY_POINT"],
		BindAddr:          m["UDOCS_BIND_ADDR"],
		Port:              m["UDOCS_PORT"],
		RootRoute:         m["UDOCS_ROOT_ROUTE"],
		Organization:      m["UDOCS_ORGANIZATION"],
		Email:             m["UDOCS_EMAIL"],
		SearchPlaceholder: m["UDOCS_SEARCH_PLACEHOLDER"],
		Routes:            strings.Split(m["UDOCS_ROUTES"], ","),
		Seed:              strings.Split(m["UDOCS_SEED"], ","),
		MongoURL:          m["UDOCS_MONGO_URL"],
		QuipAccessToken:   m["UDOCS_QUIP_ACCESS_TOKEN"],
	}
}

func toMap(r io.Reader) map[string]string {
	conf := make(map[string]string)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		tokens := strings.Split(scanner.Text(), "=")
		if len(tokens) == 2 {
			conf[tokens[0]] = tokens[1]
		}
	}
	return conf
}

func sliceToString(slice []string) string {
	buf := new(bytes.Buffer)
	for i, v := range slice {
		buf.WriteString(v)
		if i != len(slice)-1 {
			buf.WriteRune(44) // comma
		}
	}
	return buf.String()
}
