package controller

import (
	"payment-system/internal/db/mysql"
	"payment-system/internal/restapi/models"
)

func (c *Controller) CreateWallet(cwr models.CreateWalletRequest) (uint, error) {
	return c.mysqlDB.AddWallet(mysql.AddWalletTask{
		Name:        cwr.Name,
		Description: cwr.Description,
	})
}

func (c *Controller) Wallet(id uint) (models.WalletResponse, error) {
	wallet, err := c.mysqlDB.GetWallet(id)
	if err != nil {
		return models.WalletResponse{}, err
	}
	return models.WalletResponse{
		Status:      models.OK,
		Id:          wallet.ID,
		Name:        wallet.Name,
		Description: wallet.Description,
		Balance:     wallet.Balance,
	}, nil
}
