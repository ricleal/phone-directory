START TRANSACTION;

ALTER TABLE
    `users`
ADD
    `phone` VARCHAR(15) NOT NULL;

-- copy existing phone numbers from phones table to users table
-- assuming 1 phone number per user
UPDATE
    users,
    phones
SET
    users.phone = phones.phone
WHERE
    phones.deleted_at IS NULL
    AND users.deleted_at IS NULL
    AND users.id = phones.user_id;

DROP TABLE IF EXISTS `phones`;

COMMIT;