package main

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strings"
)

type RemoteAddress struct {
	Address string `json:"address"`
	Port    string `json:"port"`
}

type Extras struct {
	Envs    map[string]string `json:"envs,omitempty"`
	AppName string            `json:"app_name,omitempty"`
}

type Response struct {
	Headers     map[string][]string `json:"headers"`
	Form        map[string][]string `json:"form"`
	Query       map[string][]string `json:"query"`
	Remote      RemoteAddress       `json:"remote"`
	Path        string              `json:"path"`
	Method      string              `json:"method"`
	ContentType string              `json:"content-type,omitempty"`
	Extras      Extras              `json:"extras,omitempty"`
	Json        json.RawMessage     `json:"json,omitempty"`
}

func parseRemoteAddr(remoteAddress string) RemoteAddress {
	lastIndex := strings.LastIndex(remoteAddress, ":")
	address := remoteAddress[:lastIndex]
	port := remoteAddress[lastIndex+1:]

	return RemoteAddress{
		Address: address,
		Port:    port,
	}
}

func echo(extras Extras) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		contentType := r.Header.Get("Content-Type")
		var jsonBody json.RawMessage = []byte("")

		if r.Body != nil {
			defer r.Body.Close()
			if strings.Contains(contentType, "application/json") {
				body, err := io.ReadAll(r.Body)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				jsonBody = json.RawMessage(body)
			}
		}
		response := Response{
			Path:        r.URL.Path,
			Method:      r.Method,
			Headers:     r.Header,
			Form:        r.Form,
			Query:       r.URL.Query(),
			Remote:      parseRemoteAddr(r.RemoteAddr),
			ContentType: r.Header.Get("Content-Type"),
			Json:        jsonBody,
			Extras:      extras,
		}

		bytes, err := json.Marshal(response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(bytes)
		if err != nil {
			slog.Error(err.Error())
		}
		slog.Info("", "response", response)
	}
}
