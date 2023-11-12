START TRANSACTION;

CREATE TABLE `phones` (
    `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `phone` varchar(15) NOT NULL,
    `user_id` bigint(20) unsigned NOT NULL,
    `created_at` DATETIME NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME DEFAULT NULL,
    `deleted_at` DATETIME DEFAULT NULL,
    FOREIGN KEY (`user_id`) REFERENCES `users` (`id`),
    PRIMARY KEY (`id`),
    KEY `idx_companies_deleted_at` (`deleted_at`)
);

-- copy existing phone numbers from users table
INSERT INTO
    `phones` (`phone`, `user_id`)
SELECT
    `phone`,
    `id`
FROM
    `users`;

-- drop column phone from users table
ALTER TABLE
    `users` DROP COLUMN `phone`;

COMMIT;