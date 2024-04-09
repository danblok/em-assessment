CREATE TABLE IF NOT EXISTS cars (
    car_id uuid DEFAULT gen_random_uuid(),
    reg_num TEXT UNIQUE NOT NULL,
    mark TEXT NOT NULL,
    model TEXT NOT NULL,
    year INT,
    PRIMARY KEY(car_id)
);

