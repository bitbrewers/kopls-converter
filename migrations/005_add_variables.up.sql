ALTER TABLE door_models
ADD COLUMN slate_position INTEGER DEFAULT 0;

ALTER TABLE programs
ADD COLUMN slate_hinge INTEGER DEFAULT 0;

UPDATE door_models SET slate_position = 0;

UPDATE programs SET slate_hinge = 0;
