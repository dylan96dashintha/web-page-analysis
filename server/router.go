package server

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/web-page-analysis/bootstrap"
	"github.com/web-page-analysis/container"
	"github.com/web-page-analysis/server/endpoint"
	"github.com/web-page-analysis/server/middleware"
	"net/http"
	"time"
)

func InitRouter(ctx context.Context, conf bootstrap.Config, ctr container.Container) {
	r := mux.NewRouter()

	analyserObj := endpoint.NewAnalyser(ctr)
	r.HandleFunc("/analyse", analyserObj.Analyse).Methods(http.MethodPost, http.MethodOptions)
	corsHandler := middleware.CorsMiddleware(r)
	server := &http.Server{
		Addr:         fmt.Sprintf("%v:%v", "0.0.0.0", conf.AppConfig.Port),
		WriteTimeout: time.Second * 350,
		ReadTimeout:  time.Second * 350,
		IdleTimeout:  time.Second * 600,
		Handler:      corsHandler,
	}
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.WithContext(ctx).Fatalf("http server error: %+v", err)
		}
	}()
}
