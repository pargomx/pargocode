INSERT INTO migraciones VALUES (02, CURRENT_TIMESTAMP, "Agregar zero_is_null a campos");

ALTER TABLE campos ADD COLUMN zero_is_null INT NOT NULL DEFAULT 0;