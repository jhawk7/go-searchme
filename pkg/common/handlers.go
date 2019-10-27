package common

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func ErrorHandler(w http.ResponseWriter, err error) {
	if err != nil {
		fmt.Println(err)
		json.NewEncoder(w).Encode(err)
	}
}