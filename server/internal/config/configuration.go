package config

import (
	"encoding/json"
	"fmt"
	"strings"

	"golang.org/x/oauth2"
)

type Mode uint8

const (
	InMemory Mode = iota + 1
	Persistent
	Cloud
)

var (
	Mode_name = map[uint8]string{
		1: "inmemory",
		2: "persistent",
		3: "cloud",
	}
	Mode_value = map[string]uint8{
		"inmemory":   1,
		"persistent": 2,
		"cloud":      3,
	}
)

// String allows Mode to implement fmt.Stringer
func (s Mode) String() string {
	return Mode_name[uint8(s)]
}

func (s *Mode) UnmarshalJSON(data []byte) (err error) {
	var modes string
	if err := json.Unmarshal(data, &modes); err != nil {
		return err
	}
	if *s, err = ParseMode(modes); err != nil {
		return err
	}
	return nil
}

// Convert a string to a Mode, returns an error if the string is unknown.
// NOTE: for JSON marshaling this must return a Mode value not a pointer, which is
// common when using integer enumerations (or any primitive type alias).
func ParseMode(s string) (Mode, error) {
	s = strings.TrimSpace(strings.ToLower(s))
	value, ok := Mode_value[s]
	if !ok {
		return Mode(0), fmt.Errorf("%q is not a valid mode", s)
	}
	return Mode(value), nil
}

type ApiServer struct {
	ApplicationUrl string
}

type TweetsStorage struct {
	ConnectionString string
	DatabaseName     string
}

type FeedsStorage struct {
	ConnectionString string
	DatabaseName     string
	CollectionName   string
}

type Authentication struct {
	Enable bool
	OAuth2 oauth2.Config
}

type Configuration struct {
	Mode           Mode
	ApiServer      ApiServer
	TweetsStorage  TweetsStorage
	FeedsStorage   FeedsStorage
	NATSUrl        string
	Authentication Authentication
	RedirectURI    string
	AllowOrigin    string // When using credentials (like cookies or HTTP authentication), CORS header Access-Control-Allow-Origin cannot be set to *
	ProjectId      string
}
