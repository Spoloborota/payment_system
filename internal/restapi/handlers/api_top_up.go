package handlers

import (
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
	"payment-system/internal/controller"
	"payment-system/internal/restapi/models"
)

type TopUpWrapper struct {
	Logger *zap.Logger
	Cntrlr *controller.Controller
}

func (x *TopUpWrapper) TopUp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var twr models.TopupWalletRequest
	err := json.NewDecoder(r.Body).Decode(&twr)
	if err != nil {
		x.Logger.Error("Failed to unmarshall request", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write(errResp(err))
		return
	}
	//TODO check whether amount > 0
	err = x.Cntrlr.TopupWallet(twr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write(errResp(err))
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(topupWalletResponse())
}

func topupWalletResponse() []byte {
	res := models.TopupWalletResponse{
		Status:      models.OK,
		Description: "Top Up Successful",
	}
	marshalled, _ := json.Marshal(res)
	return marshalled
}
