package extensions

import (
	"net/http"
	"twitter-clone/internal/config"
)

func EnableCors(w *http.ResponseWriter, config config.Configuration) {
	(*w).Header().Set("Access-Control-Allow-Origin", config.AllowOrigin)
	(*w).Header().Set("Access-Control-Allow-Credentials", "true")
	(*w).Header().Set("Access-Control-Allow-Headers", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "*")
}
