CREATE TABLE IF NOT EXISTS leagues (
    id INT PRIMARY KEY,
    sport_id INT,
    country_id INT,
    name TEXT,
    active BOOLEAN,
    short_code TEXT,
    image_path TEXT,
    type TEXT,
    sub_type TEXT
); 
