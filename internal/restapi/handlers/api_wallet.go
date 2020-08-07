package handlers

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
	"payment-system/internal/controller"
	"payment-system/internal/restapi/models"
	"strconv"
)

type CreateWalletWrapper struct {
	Logger *zap.Logger
	Cntrlr *controller.Controller
}

func (x *CreateWalletWrapper) CreateWallet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var cwr models.CreateWalletRequest
	err := json.NewDecoder(r.Body).Decode(&cwr)
	if err != nil {
		x.Logger.Error("Failed to unmarshall request", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write(errResp(err))
		return
	}

	walletID, err := x.Cntrlr.CreateWallet(cwr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write(errResp(err))
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(createWalletResponse(walletID))
}

type WalletWrapper struct {
	Logger *zap.Logger
	Cntrlr *controller.Controller
}

func (x *WalletWrapper) Wallet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	params := mux.Vars(r)
	idStr, ok := params["id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write(errResp(errors.New("ID param not found")))
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write(errResp(err))
		return
	}
	wallet, err := x.Cntrlr.Wallet(uint(id))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write(errResp(err))
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(walletResponse(wallet))
}

func errResp(err error) []byte {
	errorResp := models.ErrorResponse{Errors: &models.ErrorResponseErrors{
		Status: string(models.ERROR_),
		Body:   err.Error()}}
	marshalled, _ := json.Marshal(errorResp)
	return marshalled
}

func createWalletResponse(walletID uint) []byte {
	res := models.CreateWalletResponse{
		Status:          models.OK,
		CreatedWalletId: walletID,
	}
	marshalled, _ := json.Marshal(res)
	return marshalled
}

func walletResponse(wallet models.WalletResponse) []byte {
	marshalled, _ := json.Marshal(wallet)
	return marshalled
}
