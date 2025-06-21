package server

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/web-page-analysis/bootstrap"
	"github.com/web-page-analysis/container"
	"github.com/web-page-analysis/server/endpoint"
	"net/http"
	"time"
)

func InitRouter(ctx context.Context, conf bootstrap.Config, ctr container.Container) {
	r := mux.NewRouter()

	analyserObj := endpoint.NewAnalyser(ctr)
	r.HandleFunc("/analyse", analyserObj.Analyse).Methods(http.MethodPost)

	server := &http.Server{
		Addr: fmt.Sprintf("%v:%v", "0.0.0.0", conf.AppConfig.Port),
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 350,
		ReadTimeout:  time.Second * 350,
		IdleTimeout:  time.Second * 600,
		Handler:      r,
	}
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.WithContext(ctx).Fatalf("http server error: %+v", err)
		}
	}()
}
