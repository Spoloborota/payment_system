package controller

import (
	"fmt"
	"payment-system/internal/db"
	"payment-system/internal/restapi/models"
)

func (c *Controller) TransactionStatus(id uint64) (models.TransactionStatusResponse, error) {
	status, statusDescription, err := c.mysqlDB.GetTransactionStatus(id)
	if err != nil {
		return models.TransactionStatusResponse{}, err
	}
	var statusResult string
	switch status {
	case db.Registered:
		statusResult = "Transaction registered."
	case db.Failed:
		statusResult = "Transaction failed."
	case db.Success:
		statusResult = "Transaction succeeded."
		//TODO default case
	}
	if statusDescription != "" {
		statusResult = fmt.Sprintf("%s Description: %s.", statusResult, statusDescription)
	}
	return models.TransactionStatusResponse{
		Status:      models.OK,
		Description: statusResult,
	}, nil
}
