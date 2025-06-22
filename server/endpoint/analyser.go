package endpoint

import (
	"context"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/web-page-analysis/container"
	"github.com/web-page-analysis/domain"
	erro "github.com/web-page-analysis/server/error"
	"github.com/web-page-analysis/service"
	"net/http"
)

type Analyser struct {
	container container.Container
}

func NewAnalyser(ctr container.Container) *Analyser {
	return &Analyser{
		container: ctr,
	}
}

func (a Analyser) Analyse(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	log.WithContext(ctx).Info("start to analyse the web pages")

	// unmarshal the request
	var analyserRequest domain.AnalyserRequest
	err := json.NewDecoder(r.Body).Decode(&analyserRequest)
	if err != nil {
		log.Errorf("ERROR decoding request body, err: %+v", err)
		erro.BadRequestError(fmt.Sprintf("ERROR decoding request body, err: %+v",
			err), w)
		return
	}
	analyser := service.NewAnalyser(a.container)
	result, err := analyser.WebAnalyser(ctx, analyserRequest)
	if err != nil {
		erro.GeneralError(fmt.Sprintf("err: %+v",
			err), "error in analysing the webpage", http.StatusInternalServerError, w)
		return
	}
	raw, err := json.Marshal(result)
	if err != nil {
		erro.GeneralError(fmt.Sprintf("err: %+v",
			err), "error in marshalling response", http.StatusInternalServerError, w)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(raw)
	return
}
