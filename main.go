package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"go.uber.org/automaxprocs/maxprocs"
)

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

	server := &http.Server{Addr: "0.0.0.0:" + port, Handler: mux}
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
