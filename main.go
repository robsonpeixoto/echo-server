package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	_ "go.uber.org/automaxprocs"
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

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	extras := Extras{
		AppName: os.Getenv("APP_NAME"),
	}

	showEnvs := os.Getenv("SHOW_ENVS")
	if showEnvs == "1" {
		extras.Envs = map[string]string{}
		for _, kv := range os.Environ() {
			kvSlice := strings.SplitN(kv, "=", 2)
			k := kvSlice[0]
			v := kvSlice[1]
			extras.Envs[k] = v
		}
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", echo(extras))

	server := &http.Server{Addr: "0.0.0.0:5000", Handler: mux}
	serverCtx, serverStopCtx := context.WithCancel(context.Background())
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig

		shutdownCtx, shutdownCtxCancel := context.WithTimeout(serverCtx, 30*time.Second)
		defer shutdownCtxCancel()

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()

		err := server.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal(err)
		}
		serverStopCtx()
	}()

	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

	<-serverCtx.Done()
}
