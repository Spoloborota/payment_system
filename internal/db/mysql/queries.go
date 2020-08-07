package mysql

const (
	selectAllRegistered = `
SELECT id, created, debit_wallet_id, credit_wallet_id, amount, description, status
FROM payment_system.transactions
WHERE status = 0
ORDER BY created ASC;`

	insertWallet = `
INSERT INTO payment_system.wallets
(name, description, wallet_type_id)
VALUES(?, ?, 2);`

	getWallet = `
SELECT id, name, description, balance
FROM payment_system.wallets
WHERE id = ?;`

	insertTransaction = `
INSERT INTO payment_system.transactions
(debit_wallet_id, credit_wallet_id, amount, description)
VALUES(?, ?, ?, ?);`

	getTransactionStatus = `
SELECT status, status_description
FROM payment_system.transactions
WHERE id = ?;`

	updateTransactionStatus = `
UPDATE payment_system.transactions
SET status=?
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
