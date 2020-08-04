package mysql

const (
	selectAllRegistered = `
SELECT id, created, debit_wallet_id, credit_wallet_id, Amount, Description, Status
FROM payment_system.transactions
WHERE Status = 0
ORDER BY created ASC;`

	insertWallet = `
INSERT INTO payment_system.wallets
(Name, Description, wallet_type_id)
VALUES(?, ?, 2);`

	insertTransaction = `
INSERT INTO payment_system.transactions
(debit_wallet_id, credit_wallet_id, Amount, Description)
VALUES(?, ?, ?, ?);`

	updateTransactionStatus = `
UPDATE payment_system.transactions
SET Status=?
WHERE id=?;`

	topUpWallet = `
UPDATE payment_system.wallets
SET balance = balance + ? 
WHERE id=?;`

	getWalletBalance = `
SELECT balance
FROM payment_system.wallets
WHERE id = ?;`

	updateDebitWalletBalance = `
UPDATE payment_system.wallets
SET balance = balance + ?
WHERE id=?;`

	updateCreditWalletBalance = `
UPDATE payment_system.wallets
SET balance = balance - ?
WHERE id=?;`
)
