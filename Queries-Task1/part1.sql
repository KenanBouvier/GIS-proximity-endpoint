UPDATE "MY_TABLE"
SET website = substring(website from '(?:.*://)?(?:www\.)?([^/?]*)');

-- Applying a regex pattern given the website field and updating its value across all records.
