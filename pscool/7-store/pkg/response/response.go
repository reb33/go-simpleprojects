package response

import (
	"encoding/json"
	"net/http"
)



func Json(w http.ResponseWriter, code int, res any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(res)
}

func PlainText(w http.ResponseWriter, code int, res string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(code)
	w.Write([]byte(res))
}