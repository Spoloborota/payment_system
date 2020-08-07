package handlers

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"payment-system/internal/controller"
	"payment-system/internal/restapi/models"
)

type CreateTransactionWrapper struct {
	Logger *zap.Logger
	Cntrlr *controller.Controller
}

func (x *CreateTransactionWrapper) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var ctr models.CreateTransactionRequest
	err := json.NewDecoder(r.Body).Decode(&ctr)
	if err != nil {
		x.Logger.Error("Failed to unmarshall request", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write(errResp(err))
		return
	}
	//TODO check whether amount > 0
	id, err := x.Cntrlr.AddTransaction(ctr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write(errResp(err))
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(createTransactionResponse(id))
}

func createTransactionResponse(id uint64) []byte {
	res := models.CreateTransactionResponse{
		Status:      models.OK,
		Description: fmt.Sprintf("Transaction registered: %d", id),
	}
	marshalled, _ := json.Marshal(res)
	return marshalled
}
