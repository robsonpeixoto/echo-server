package echo

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
	Host          string              `json:"host"`
	Proto         string              `json:"proto"`
	ContentLength int64               `json:"content_length"`
	Headers       map[string][]string `json:"headers"`
	Form          map[string][]string `json:"form"`
	Query         map[string][]string `json:"query"`
	Remote        RemoteAddress       `json:"remote"`
	Path          string              `json:"path"`
	Method        string              `json:"method"`
	ContentType   string              `json:"content-type,omitempty"`
	Extras        Extras              `json:"extras"`
	JSON          json.RawMessage     `json:"json,omitempty"`
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

func Handler(extras Extras) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		contentType := r.Header.Get("Content-Type")
		var jsonBody json.RawMessage = []byte("")

		if r.Body != nil {
			defer func() {
				if err := r.Body.Close(); err != nil {
					slog.Error("failed to close request body", "error", err)
				}
			}()
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
			Host:          r.Host,
			Proto:         r.Proto,
			ContentLength: r.ContentLength,
			Path:          r.URL.Path,
			Method:        r.Method,
			Headers:       r.Header,
			Form:          r.Form,
			Query:         r.URL.Query(),
			Remote:        parseRemoteAddr(r.RemoteAddr),
			ContentType:   r.Header.Get("Content-Type"),
			JSON:          jsonBody,
			Extras:        extras,
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
