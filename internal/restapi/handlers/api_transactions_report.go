package handlers

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"go.uber.org/zap"
	"net/http"
	"payment-system/internal/controller"
	"payment-system/internal/restapi/models"
	"strconv"
)

type TransactionsReportWrapper struct {
	Logger *zap.Logger
	Cntrlr *controller.Controller
}

func (x *TransactionsReportWrapper) TransactionsReport(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/csv; charset=UTF-8")
	var filter models.TransactionReportRequest
	isTopUpStr, ok := r.URL.Query()["is_top_up"]
	if ok {
		parseBool, err := strconv.ParseBool(isTopUpStr[0])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write(errResp(err))
			return
		}
		filter.IsTopUps = parseBool
	}
	startDateStr, ok := r.URL.Query()["start_date"]
	if ok {
		startDate, err := strconv.ParseInt(startDateStr[0], 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write(errResp(err))
			return
		}
		filter.CreatedFrom = startDate
	}
	endDateStr, ok := r.URL.Query()["end_date"]
	if ok {
		endDate, err := strconv.ParseInt(endDateStr[0], 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write(errResp(err))
			return
		}
		filter.CreatedTo = endDate
	}
	debitWalletIDStr, ok := r.URL.Query()["debit_wallet_id"]
	if ok {
		debitWalletID, err := strconv.ParseUint(debitWalletIDStr[0], 10, 32)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write(errResp(err))
			return
		}
		filter.DebitWalletID = uint32(debitWalletID)
	}
	creditWalletIDStr, ok := r.URL.Query()["credit_wallet_id"]
	if ok {
		creditWalletID, err := strconv.ParseUint(creditWalletIDStr[0], 10, 32)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write(errResp(err))
			return
		}
		filter.CreditWalletID = uint32(creditWalletID)
	}

	transactions, err := x.Cntrlr.TransactionsReport(filter)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write(errResp(err))
		return
	}

	var buffer bytes.Buffer
	bufferWriter := bufio.NewWriter(&buffer)
	writer := csv.NewWriter(bufferWriter)
	writer.Comma = ';'

	err = writer.WriteAll(transactions)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write(errResp(err))
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(buffer.Bytes())
}
