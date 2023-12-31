package main

import (
	_ "fmt"
	"twitter-clone/internal/config"
)

func main() {
	configuration := config.ReadConfiguration()
	_ = configuration
}
