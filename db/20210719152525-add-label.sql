-- Add label table and foreign key in comments

CREATE TABLE IF NOT EXISTS labels (
  labelHex               TEXT          NOT NULL  UNIQUE  PRIMARY KEY      ,
  color                  TEXT          NOT NULL                           ,
  domain                 TEXT          NOT NULL
);

ALTER TABLE comments ADD COLUMN IF NOT EXISTS labelHex TEXT;
