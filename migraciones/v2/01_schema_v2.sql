INSERT INTO migraciones VALUES (2,1, CURRENT_TIMESTAMP, 'Esquema completo v2');

CREATE TABLE paquetes (
  paquete_id INT NOT NULL,
  go_module TEXT NOT NULL DEFAULT '',
  directorio TEXT NOT NULL DEFAULT '',
  nombre TEXT NOT NULL,
  descripcion TEXT NOT NULL DEFAULT '',
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
  humano_plural TEXT NOT NULL DEFAULT '',
  kebab TEXT NOT NULL,
  es_femenino INT NOT NULL,
  descripcion TEXT NOT NULL DEFAULT '',
  directrices TEXT NOT NULL DEFAULT '',
  PRIMARY KEY (tabla_id),
  UNIQUE (paquete_id, nombre_repo),
  UNIQUE (paquete_id, abrev),
  FOREIGN KEY (paquete_id) REFERENCES paquetes (paquete_id) ON DELETE RESTRICT ON UPDATE CASCADE
);

CREATE TABLE campos (
  campo_id INT NOT NULL,
  tabla_id INT NOT NULL,
  nombre_campo TEXT NOT NULL,
  nombre_columna TEXT NOT NULL,
  nombre_humano TEXT NOT NULL,
  tipo_go TEXT NOT NULL DEFAULT '',
  tipo_sql TEXT NOT NULL DEFAULT '',
  setter TEXT NOT NULL DEFAULT '',
  importado INT NOT NULL,
  primary_key INT NOT NULL,
  foreign_key INT NOT NULL,
  uq INT NOT NULL,
  req INT NOT NULL,
  ro INT NOT NULL,
  filtro INT NOT NULL,
  nullable INT NOT NULL,
  zero_is_null INT NOT NULL DEFAULT 0,
  max_lenght INT NOT NULL,
  uns INT NOT NULL,
  default_sql TEXT NOT NULL DEFAULT '',
  especial INT NOT NULL,
  referencia_campo INT,
  expresion TEXT NOT NULL DEFAULT '',
  es_femenino INT NOT NULL,
  descripcion TEXT NOT NULL DEFAULT '',
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
  descripcion TEXT NOT NULL DEFAULT '',
  PRIMARY KEY (campo_id,clave),
  FOREIGN KEY (campo_id) REFERENCES campos (campo_id) ON DELETE RESTRICT ON UPDATE CASCADE
);

CREATE TABLE consultas (
  consulta_id INT NOT NULL,
  paquete_id INT NOT NULL,
  tabla_id INT NOT NULL,
  nombre_item TEXT NOT NULL,
  nombre_items TEXT NOT NULL DEFAULT '',
  abrev TEXT NOT NULL DEFAULT '',
  es_femenino INT NOT NULL,
  descripcion TEXT NOT NULL DEFAULT '',
  directrices TEXT NOT NULL DEFAULT '',
  PRIMARY KEY (consulta_id),
  FOREIGN KEY (paquete_id) REFERENCES paquetes (paquete_id) ON DELETE RESTRICT ON UPDATE CASCADE,
  FOREIGN KEY (tabla_id) REFERENCES tablas (tabla_id) ON DELETE RESTRICT ON UPDATE CASCADE
);

CREATE TABLE consulta_relaciones (
  consulta_id INT NOT NULL,
  posicion INT NOT NULL,
  tipo_join NOT NULL DEFAULT '',
  join_as TEXT NOT NULL DEFAULT '',
  join_tabla_id INT NOT NULL,
  join_on TEXT NOT NULL DEFAULT '',
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
  expresion TEXT NOT NULL DEFAULT '',
  alias_sql TEXT NOT NULL DEFAULT '',
  nombre_campo TEXT NOT NULL,
  tipo_go TEXT NOT NULL,
  pk INT NOT NULL,
  filtro INT NOT NULL,
  group_by INT NOT NULL,
  descripcion TEXT NOT NULL DEFAULT '',
  PRIMARY KEY (consulta_id,posicion),
  FOREIGN KEY (consulta_id) REFERENCES consultas (consulta_id) ON DELETE RESTRICT ON UPDATE CASCADE,
  FOREIGN KEY (campo_id) REFERENCES campos (campo_id) ON DELETE RESTRICT ON UPDATE CASCADE
);