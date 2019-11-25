package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

/*
	Message
*/
func Message(status bool, message string) map[string]interface{} {
	return map[string]interface{}{"status": status, "message": message}
}

func MessageWithData(status bool, message string, data interface{}) map[string]interface{} {
	return map[string]interface{}{"status": status, "message": message, "data": data}
}

func Respond(w http.ResponseWriter, data map[string]interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "content-type,Authorization")
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE")
	if json.NewEncoder(w).Encode(data) != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func GetRequestParam(r *http.Request, key string) string {
	vals := r.URL.Query()[key]

	fmt.Println("GetRequestParam, vals: ", vals)

	if vals == nil || len(vals) == 0 {
		return ""
	}

	return vals[0]
}
