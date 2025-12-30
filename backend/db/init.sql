PRAGMA foreign_keys=ON;
CREATE TABLE IF NOT EXISTS project (
   ID INTEGER PRIMARY KEY,
   name VARCHAR(255),
   route TEXT, -- con esto puedo saber el entorno para ejecutar
   command TEXT -- con esto lo empiezo a correr
);

CREATE TABLE IF NOT EXISTS execution_logs(
   ID INTEGER PRIMARY KEY,
   id_project
   content TEXT
);