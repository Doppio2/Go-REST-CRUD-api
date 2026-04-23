-- +goose Up
CREATE TABLE IF NOT EXISTS equipment (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    description TEXT,
    creation_date TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS experiment (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    description TEXT,
    creation_date TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS experiment_equipment (
    experiment_id INTEGER NOT NULL,
    equipment_id INTEGER NOT NULL,
    PRIMARY KEY (experiment_id, equipment_id),
    FOREIGN KEY (experiment_id) REFERENCES experiment(id) ON DELETE CASCADE,
    FOREIGN KEY (equipment_id) REFERENCES equipment(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE IF EXISTS experiment_equipment;
DROP TABLE IF EXISTS experiment;
DROP TABLE IF EXISTS equipment;
