/* Hay configuraciones PRAGMA que aplican para cada conexión y no deben estar aquí */
PRAGMA journal_mode = WAL;

/* Crear tabla para control de migraciones */
CREATE TABLE migraciones (
	id INT NOT NULL,
	fecha TEXT NOT NULL,
	detalles TEXT NOT NULL,
	PRIMARY KEY (id),
	UNIQUE(detalles)
);

/* Migración origen */
INSERT INTO migraciones VALUES (00, CURRENT_TIMESTAMP, "Crear tabla para migraciones");
