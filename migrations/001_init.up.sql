CREATE TABLE IF NOT EXISTS people (
        id SERIAL PRIMARY KEY,
        name TEXT NOT NULL,
        surname TEXT NOT NULL,
        patronymic TEXT,
        gender TEXT,
        age INT,
        nationality TEXT
);

CREATE UNIQUE INDEX idx_people_name ON people(name, surname);


