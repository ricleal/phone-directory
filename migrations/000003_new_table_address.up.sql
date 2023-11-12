CREATE table addresses (
    id bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    address varchar(255) NOT NULL,
    user_id bigint(20) unsigned NOT NULL,
    created_at DATETIME NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT NULL,
    deleted_at DATETIME DEFAULT NULL,
    FOREIGN KEY (`user_id`) REFERENCES `users` (`id`),
    PRIMARY KEY (`id`),
    KEY `idx_companies_deleted_at` (`deleted_at`)
);
