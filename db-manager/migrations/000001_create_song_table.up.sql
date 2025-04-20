CREATE TABLE IF NOT EXISTS SONGS (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    audio TEXT NOT NULL,
    bg TEXT NOT NULL,
    osu_file TEXT NOT NULL,
    path TEXT NOT NULL,
    date_added TEXT NOT NULL,
    title TEXT NOT NULL,
    artist TEXT NOT NULL,
    creator TEXT NOT NULL,
    duration REAL NOT NULL,
    beatmap_set_id INTEGER,
    mode INTEGER,
    title_unicode TEXT,
    artist_unicode TEXT,
    primary_color TEXT,
    secondary_color TEXT
);