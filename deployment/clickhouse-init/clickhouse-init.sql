CREATE TABLE IF NOT EXISTS `transactions` (
    id UInt32,
    created UInt64,
    debit_wallet_id UInt32,
    credit_wallet_id UInt32,
    amount UInt32
)
engine = MergeTree() PARTITION BY toYYYYMM(toDate(created)) ORDER BY created SETTINGS index_granularity = 8192;