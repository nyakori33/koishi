package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

var (
	client http.Client

	Port = flag.Uint("port", 3020, "Listen TCP port.")
)

func init() {
	flag.Parse()
}

func main() {
	http.HandleFunc("/", handler)
	addr := fmt.Sprintf(":%d", *Port)
	error := http.ListenAndServe(addr, nil)
	if error != nil {
		panic(error)
	}
}

func handler(html http.ResponseWriter, remote *http.Request) {
	uri := remote.FormValue("uri")
	uri, error := url.QueryUnescape(uri)
	if error != nil {
		html.WriteHeader(http.StatusInternalServerError)
		html.Write([]byte(error.Error()))
		return
	}

	data, error := io.ReadAll(remote.Body)
	if error != nil {
		html.WriteHeader(http.StatusInternalServerError)
		html.Write([]byte(error.Error()))
		return
	}

	request, error := http.NewRequest(remote.Method, uri, bytes.NewBuffer(data))
	if error != nil {
		html.WriteHeader(http.StatusInternalServerError)
		html.Write([]byte(error.Error()))
		return
	}

	for key, data := range remote.Header {
		value := strings.Join(data, " ")
		request.Header.Add(key, value)
	}

	response, error := client.Do(request)
	if error != nil {
		html.WriteHeader(500)
		html.Write([]byte(error.Error()))
		return
	}

	defer response.Body.Close()

	for key, value := range response.Header {
		html.Header().Add(key, strings.Join(value, " "))
	}

	data, error = io.ReadAll(response.Body)
	if error != nil {
		html.WriteHeader(500)
		html.Write([]byte(error.Error()))
		return
	}

	html.Write(data)
}
