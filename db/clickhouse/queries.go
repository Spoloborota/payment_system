package clickhouse

const (
	insertTransaction = `
INSERT INTO payment_system.transactions
(id, created, debit_wallet_id, credit_wallet_id, amount, description, status)
VALUES(?, ?, ?, ?, ?, ?, ?);`
)
