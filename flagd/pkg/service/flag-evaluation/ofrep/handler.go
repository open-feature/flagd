package ofrep

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/open-feature/flagd/core/pkg/evaluator"
	"github.com/open-feature/flagd/core/pkg/logger"
)

type handler struct {
	Logger    *logger.Logger
	evaluator evaluator.IEvaluator
}

func NewHandler(logger *logger.Logger, evaluator evaluator.IEvaluator) http.Handler {
	h := handler{
		logger,
		evaluator,
	}

	router := mux.NewRouter()
	router.HandleFunc("/ofrep/v1/evaluate/flags/{key}", h.HandleFlagEvaluation)
	router.HandleFunc("/ofrep/v1/evaluate/flags", h.HandleBulkEvaluation)

	return router
}

func (h *handler) HandleFlagEvaluation(_ http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Println(vars["key"])

	//  need new evaluator mechanism to evaluate unknown type

	h.Logger.Info("single flag")
}

func (h *handler) HandleBulkEvaluation(_ http.ResponseWriter, _ *http.Request) {
	h.Logger.Info("bulk flag")
}
