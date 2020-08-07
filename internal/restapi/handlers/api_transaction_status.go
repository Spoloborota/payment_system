package handlers

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
	"payment-system/internal/controller"
	"strconv"
)

type TransactionStatusWrapper struct {
	Logger *zap.Logger
	Cntrlr *controller.Controller
}

func (x *TransactionStatusWrapper) TransactionStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	params := mux.Vars(r)
	id, ok := params["id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write(errResp(errors.New("ID param not found")))
		return
	}

	parsedID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		_, _ = w.Write(errResp(errors.New("incorrect id")))
		return
	}
	status, err := x.Cntrlr.TransactionStatus(uint64(parsedID))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write(errResp(err))
		return
	}
	marshalled, _ := json.Marshal(status)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(marshalled)
}
