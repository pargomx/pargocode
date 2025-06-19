INSERT INTO migraciones VALUES (01, CURRENT_TIMESTAMP, "Esquema: domain driven design");

CREATE TABLE paquetes (
  paquete_id INT NOT NULL,
  go_module TEXT NOT NULL,
  directorio TEXT NOT NULL,
  nombre TEXT NOT NULL,
  descripcion TEXT NOT NULL,
  PRIMARY KEY (paquete_id)
);

CREATE TABLE tablas (
  tabla_id INT NOT NULL,
  paquete_id INT NOT NULL,
  nombre_repo TEXT NOT NULL,
  nombre_item TEXT NOT NULL,
  nombre_items TEXT NOT NULL,
  abrev TEXT NOT NULL,
  humano TEXT NOT NULL,
  humano_plural TEXT NOT NULL,
  kebab TEXT NOT NULL,
  es_femenino NOT NULL,
  descripcion TEXT NOT NULL,
  directrices TEXT NOT NULL,
  PRIMARY KEY (tabla_id),
  UNIQUE (nombre_repo),
  UNIQUE (abrev),
  FOREIGN KEY (paquete_id) REFERENCES paquetes (paquete_id) ON DELETE RESTRICT ON UPDATE CASCADE
);

CREATE TABLE campos (
  campo_id INT NOT NULL,
  tabla_id INT NOT NULL,
  nombre_campo TEXT NOT NULL,
  nombre_columna TEXT NOT NULL,
  nombre_humano TEXT NOT NULL,
  tipo_go TEXT NOT NULL,
  tipo_sql TEXT NOT NULL,
  setter TEXT NOT NULL,
  importado NOT NULL,
  primary_key NOT NULL,
  foreign_key NOT NULL,
  uq NOT NULL,
  req NOT NULL,
  ro NOT NULL,
  filtro NOT NULL,
  nullable NOT NULL,
  max_lenght INT NOT NULL,
  uns NOT NULL,
  default_sql TEXT NOT NULL,
  especial NOT NULL,
  referencia_campo INT,
  expresion TEXT NOT NULL,
  es_femenino NOT NULL,
  descripcion TEXT NOT NULL,
  posicion INT NOT NULL,
  PRIMARY KEY (campo_id),
  UNIQUE (tabla_id, nombre_campo),
  UNIQUE (tabla_id, nombre_columna),
  FOREIGN KEY (tabla_id) REFERENCES tablas (tabla_id) ON DELETE RESTRICT ON UPDATE CASCADE
);

CREATE TABLE valores_enum (
  campo_id INT NOT NULL,
  numero INT NOT NULL,
  clave TEXT NOT NULL,
  etiqueta TEXT NOT NULL,
  descripcion TEXT NOT NULL,
  PRIMARY KEY (campo_id,clave),
  FOREIGN KEY (campo_id) REFERENCES campos (campo_id) ON DELETE RESTRICT ON UPDATE CASCADE
);

/* Consultas */

CREATE TABLE consultas (
  consulta_id INT NOT NULL,
  paquete_id INT NOT NULL,
  tabla_id INT NOT NULL,
  nombre_item TEXT NOT NULL,
  nombre_items TEXT NOT NULL,
  abrev TEXT NOT NULL,
  es_femenino NOT NULL,
  descripcion TEXT NOT NULL,
  directrices TEXT NOT NULL,
  PRIMARY KEY (consulta_id),
  FOREIGN KEY (paquete_id) REFERENCES paquetes (paquete_id) ON DELETE RESTRICT ON UPDATE CASCADE,
  FOREIGN KEY (tabla_id) REFERENCES tablas (tabla_id) ON DELETE RESTRICT ON UPDATE CASCADE
);

CREATE TABLE consulta_relaciones (
  consulta_id INT NOT NULL,
  posicion INT NOT NULL,
  tipo_join NOT NULL,
  join_tabla_id INT NOT NULL,
  join_as TEXT NOT NULL,
  join_on TEXT NOT NULL,
  from_tabla_id INT NOT NULL,
  PRIMARY KEY (consulta_id,posicion),
  FOREIGN KEY (consulta_id) REFERENCES consultas (consulta_id) ON DELETE RESTRICT ON UPDATE CASCADE,
  FOREIGN KEY (join_tabla_id) REFERENCES tablas (tabla_id) ON DELETE RESTRICT ON UPDATE CASCADE,
  FOREIGN KEY (from_tabla_id) REFERENCES tablas (tabla_id) ON DELETE RESTRICT ON UPDATE CASCADE
);

CREATE TABLE consulta_campos (
  consulta_id INT NOT NULL,
  posicion INT NOT NULL,
  campo_id INT,
  expresion TEXT NOT NULL,
  alias_sql TEXT NOT NULL,
  nombre_campo TEXT NOT NULL,
  tipo_go TEXT NOT NULL,
  pk NOT NULL,
  filtro NOT NULL,
  group_by NOT NULL,
  descripcion TEXT NOT NULL,
  PRIMARY KEY (consulta_id,posicion),
  FOREIGN KEY (consulta_id) REFERENCES consultas (consulta_id) ON DELETE RESTRICT ON UPDATE CASCADE,
  FOREIGN KEY (campo_id) REFERENCES campos (campo_id) ON DELETE RESTRICT ON UPDATE CASCADE
);