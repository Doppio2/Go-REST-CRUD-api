-- +goose Up
CREATE INDEX IF NOT EXISTS idx_experiment_equipment_experiment_id
    ON experiment_equipment (experiment_id);

CREATE INDEX IF NOT EXISTS idx_experiment_equipment_equipment_id
    ON experiment_equipment (equipment_id);

-- +goose Down
DROP INDEX IF EXISTS idx_experiment_equipment_equipment_id;
DROP INDEX IF EXISTS idx_experiment_equipment_experiment_id;
