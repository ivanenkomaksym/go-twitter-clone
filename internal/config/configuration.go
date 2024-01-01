package config

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Mode uint8

const (
	InMemory Mode = iota + 1
	Persistent
)

var (
	Mode_name = map[uint8]string{
		1: "inmemory",
		2: "persistent",
	}
	Mode_value = map[string]uint8{
		"inmemory":   1,
		"persistent": 2,
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
	ConnectionString    string
	DatabaseName        string
	FeedsCollectionName string
}

type Configuration struct {
	Mode          Mode
	ApiServer     ApiServer
	TweetsStorage TweetsStorage
	FeedsStorage  FeedsStorage
}
