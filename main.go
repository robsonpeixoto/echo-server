package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"go.uber.org/automaxprocs/maxprocs"
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

func echo(extras Extras) func(w http.ResponseWriter, r *http.Request) {
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

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	_, _ = maxprocs.Set(maxprocs.Logger(log.Printf))

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

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	server := &http.Server{Addr: net.JoinHostPort("0.0.0.0", port), Handler: mux}
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
