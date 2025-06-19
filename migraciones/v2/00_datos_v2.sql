INSERT INTO migraciones VALUES (2,0, CURRENT_TIMESTAMP, 'Migrar datos de v1 a v2');

INSERT INTO main.paquetes (
  paquete_id,
  go_module,
  directorio,
  nombre,
  descripcion
) SELECT
  paquete_id,
  go_module,
  directorio,
  nombre,
  descripcion
  FROM old_schema.paquetes
;

INSERT INTO main.tablas (
  tabla_id,
  paquete_id,
  nombre_repo,
  nombre_item,
  nombre_items,
  abrev,
  humano,
  humano_plural,
  kebab,
  es_femenino,
  descripcion,
  directrices
) SELECT
  tabla_id,
  paquete_id,
  nombre_repo,
  nombre_item,
  nombre_items,
  abrev,
  humano,
  humano_plural,
  kebab,
  es_femenino,
  descripcion,
  directrices
  FROM old_schema.tablas
;

INSERT INTO main.campos (
  campo_id,
  tabla_id,
  nombre_campo,
  nombre_columna,
  nombre_humano,
  tipo_go,
  tipo_sql,
  setter,
  importado,
  primary_key,
  foreign_key,
  uq,
  req,
  ro,
  filtro,
  nullable,
  zero_is_null,
  max_lenght,
  uns,
  default_sql,
  especial,
  referencia_campo,
  expresion,
  es_femenino,
  descripcion,
  posicion
) SELECT
  campo_id,
  tabla_id,
  nombre_campo,
  nombre_columna,
  nombre_humano,
  tipo_go,
  tipo_sql,
  setter,
  importado,
  primary_key,
  foreign_key,
  uq,
  req,
  ro,
  filtro,
  nullable,
  zero_is_null,
  max_lenght,
  uns,
  default_sql,
  especial,
  referencia_campo,
  expresion,
  es_femenino,
  descripcion,
  posicion
  FROM old_schema.campos
;

INSERT INTO main.valores_enum (
  campo_id,
  numero,
  clave,
  etiqueta,
  descripcion
) SELECT
  campo_id,
  numero,
  clave,
  etiqueta,
  descripcion
  FROM old_schema.valores_enum
;

INSERT INTO main.consultas (
  consulta_id,
  paquete_id,
  tabla_id,
  nombre_item,
  nombre_items,
  abrev,
  es_femenino,
  descripcion,
  directrices
) SELECT
  consulta_id,
  paquete_id,
  tabla_id,
  nombre_item,
  nombre_items,
  abrev,
  es_femenino,
  descripcion,
  directrices
  FROM old_schema.consultas
;

INSERT INTO main.consulta_relaciones (
  consulta_id,
  posicion,
  tipo_join,
  join_as,
  join_tabla_id,
  join_on,
  from_tabla_id
) SELECT
  consulta_id,
  posicion,
  tipo_join,
  join_as,
  join_tabla_id,
  join_on,
  from_tabla_id
  FROM old_schema.consulta_relaciones
;

INSERT INTO main.consulta_campos (
  consulta_id,
  posicion,
  campo_id,
  expresion,
  alias_sql,
  nombre_campo,
  tipo_go,
  pk,
  filtro,
  group_by,
  descripcion
) SELECT
  consulta_id,
  posicion,
  campo_id,
  expresion,
  alias_sql,
  nombre_campo,
  tipo_go,
  pk,
  filtro,
  group_by,
  descripcion
  FROM old_schema.consulta_campos
;
