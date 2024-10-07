CREATE TABLE sessions (
  id         INTEGER  PRIMARY KEY,
  cluster_name TEXT NOT NULL,
  address     TEXT NOT NULL,
  port        TEXT NOT NULL,
  cmds        TEXT NOT NULL,
  created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);