CREATE DATABASE IF NOT EXISTS `payment_system`;
USE `payment_system`;

CREATE TABLE `wallet_types` (
  `id` int unsigned NOT NULL,
  `name` varchar(255) NOT NULL,
  `description` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

INSERT INTO `payment_system`.`wallet_types`
(id, name, description)
VALUES(1, 'bank', 'Bank account. Balance can be negative');

INSERT INTO `payment_system`.`wallet_types`
(id, name, description)
VALUES(2, 'customer', 'Customer account. Balance can not be negative');


CREATE TABLE `wallets` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(50) DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL,
  `wallet_type_id` int(10) unsigned DEFAULT NULL,
  `balance` int(10) unsigned NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  KEY `wallets_FK` (`wallet_type_id`),
  CONSTRAINT `wallets_FK` FOREIGN KEY (`wallet_type_id`) REFERENCES `wallet_types` (`id`) ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

INSERT INTO payment_system.wallets
(name, description, wallet_type_id)
VALUES('bank', 'Bank wallet', 1);


CREATE TABLE `transactions` (
  `id` bigint(10) unsigned NOT NULL AUTO_INCREMENT,
  `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `debit_wallet_id` int(10) unsigned NOT NULL,
  `credit_wallet_id` int(10) unsigned NOT NULL,
  `amount` int(10) unsigned NOT NULL,
  `description` varchar(255) DEFAULT NULL,
  `status` tinyint(3) unsigned NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  KEY `transactions_status_IDX` (`status`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
