package utils

import (
	"log"
	"net/http"
	"strconv"
)

func Response(resp http.ResponseWriter, status int, msg string) {
	resp.WriteHeader(status)
	resp.Header().Set("Content-type", "application/json")
	s := `{"code":"` + strconv.Itoa(status) + `",` + msg[1:]
	log.Printf(s)
	_, err := resp.Write([]byte(s))
	if err != nil {
		panic(err.Error())
	}
}
