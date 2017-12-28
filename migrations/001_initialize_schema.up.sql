CREATE TABLE programs (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE,
    program TEXT,
    hinge_position NUMERIC(5,1)
);

CREATE TABLE door_models (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE,
    depth INTEGER,
    variable INTEGER
);

CREATE TABLE hinges (
    id SERIAL PRIMARY KEY,
    barcode CHAR(1) UNIQUE,
    variable INTEGER
);

CREATE TABLE handednesses (
    id SERIAL PRIMARY KEY,
    barcode CHAR(1) UNIQUE,
    handedness TEXT
);

CREATE TABLE handles (
    id SERIAL PRIMARY KEY,
    barcode CHAR(1) UNIQUE,
    handle INTEGER
);

CREATE TABLE handle_positions (
    id SERIAL PRIMARY KEY,
    barcode CHAR(1) UNIQUE,
    handle_position TEXT
);
