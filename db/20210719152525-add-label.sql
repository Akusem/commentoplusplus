-- Add label table and foreign key in comments
CREATE TABLE IF NOT EXISTS labels (
  labelHex               TEXT          NOT NULL  UNIQUE  PRIMARY KEY      ,
  name                   TEXT          NOT NULL  UNIQUE                   ,
  color                  TEXT          NOT NULL                           ,
  domain                 TEXT          NOT NULL
);

-- Many to many relationship through comments and labels
CREATE TABLE IF NOT EXISTS comments_labels (
  commentHex             TEXT         REFERENCES comments (commentHex) ON UPDATE CASCADE,
  labelHex               TEXT         REFERENCES labels   (labelHex)   ON UPDATE CASCADE,
  CONSTRAINT comments_labels_pkey     PRIMARY KEY (commentHex, labelHex)
);

-- Create a field to make labels optionable
ALTER TABLE domains ADD COLUMN IF NOT EXISTS allowLabels BOOLEAN DEFAULT FALSE;
