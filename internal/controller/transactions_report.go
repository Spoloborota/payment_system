package controller

import (
	"payment-system/internal/db"
	"payment-system/internal/db/clickhouse"
	"payment-system/internal/restapi/models"
	"strconv"
	"time"
)

func (c *Controller) TransactionsReport(trr models.TransactionReportRequest) ([][]string, error) {
	txs, err := c.clickHouseDB.TransactionsReport(clickhouse.ReportQuery{
		IsTopUps:       trr.IsTopUps,
		CreatedFrom:    trr.CreatedFrom,
		CreatedTo:      trr.CreatedTo,
		DebitWalletID:  trr.DebitWalletID,
		CreditWalletID: trr.CreditWalletID,
	})
	if err != nil {
		return nil, err
	}

	records := [][]string{
		{"ID", "Created", "DebitWalletID", "CreditWalletID", "Amount", "Description", "Status", "StatusDescription"},
	}

	for _, tx := range txs {
		var statusResult string
		switch tx.Status {
		case uint8(db.Registered):
			statusResult = "Transaction registered."
		case uint8(db.Failed):
			statusResult = "Transaction failed."
		case uint8(db.Success):
			statusResult = "Transaction succeeded."
			//TODO default case
		}
		record := []string{
			strconv.FormatUint(tx.ID, 10),
			time.Unix(tx.Created, 0).String(),
			strconv.FormatUint(uint64(tx.DebitWalletID), 10),
			strconv.FormatUint(uint64(tx.CreditWalletID), 10),
			strconv.FormatInt(int64(tx.Amount), 10),
			tx.Description,
			statusResult,
			tx.StatusDescription,
		}
		records = append(records, record)
	}

	return records, nil
}
