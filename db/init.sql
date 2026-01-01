PRAGMA foreign_keys=ON;

CREATE TABLE IF NOT EXISTS project (
   id INTEGER PRIMARY KEY,
   name VARCHAR(255) NOT NULL DEFAULT '' UNIQUE,
   dir TEXT NOT NULL DEFAULT '', -- con esto puedo saber el entorno para ejecutar
   command TEXT NOT NULL DEFAULT '', -- con esto lo empiezo a correr
   thumbnail_url TEXT NOT NULL DEFAULT '',
   created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
   is_paused bool NOT NULL default false
);

CREATE TABLE IF NOT EXISTS execution_logs(
   id INTEGER PRIMARY KEY,
   id_project INTEGER NOT NULL REFERENCES PROJECT(id) ON DELETE CASCADE,
   created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
   content TEXT NOT NULL DEFAULT ''
);