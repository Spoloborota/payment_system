CREATE DATABASE IF NOT EXISTS `payment_system`;
USE `payment_system`;

CREATE TABLE IF NOT EXISTS `payment_system`.`transactions` (
    id UInt64,
    created Int64,
    debit_wallet_id UInt32,
    credit_wallet_id UInt32,
    amount UInt32,
    description String,
    status UInt8
)
engine = MergeTree() PARTITION BY toYYYYMM(toDate(created)) ORDER BY created SETTINGS index_granularity = 8192;