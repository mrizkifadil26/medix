CREATE TABLE IF NOT EXISTS files (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    absolute_path TEXT NOT NULL,
    relative_path TEXT,
    "filename" TEXT,
    extension TEXT,
    "type" TEXT,
    "size" TEXT,
    modtime TIMESTAMP,

    parent_type TEXT CHECK(parent_type IN ('movie','tv')),
    parent_id INTEGER,
    fingerprint TEXT NOT NULL UNIQUE,

    created_at TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS movie (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT,
    year TEXT,
    media_id INTEGER NOT NULL,
    source TEXT,
    "group" TEXT,
    "collection" TEXT,
    thumbnail_id INTEGER,
    tmdb_id INTEGER UNIQUE,

    created_at TIMESTAMP,
    deleted_at TIMESTAMP,

    FOREIGN KEY(media_id) REFERENCES files(id),
    FOREIGN KEY(thumbnail_id) REFERENCES thumbnail(id),
    FOREIGN KEY(tmdb_id) REFERENCES tmdb_movie(id)
);

CREATE TABLE IF NOT EXISTS tv (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT,
    source TEXT,
    "group" TEXT,
    thumbnail_id INTEGER,
    tmdb_id INTEGER UNIQUE,

    created_at TIMESTAMP,
    deleted_at TIMESTAMP,

    FOREIGN KEY(thumbnail_id) REFERENCES thumbnail(id),
    FOREIGN KEY(tmdb_id) REFERENCES tmdb_tv(id)
);

CREATE TABLE IF NOT EXISTS tv_season (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    tv_id INTEGER,
    season_number INTEGER,
    episode_count INTEGER,
    first_air_date TEXT,
    last_air_date TEXT,
    tmdb_id INTEGER,

    FOREIGN KEY(tv_id) REFERENCES tv(id)
);

CREATE TABLE IF NOT EXISTS tv_episode (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    season_id INTEGER,
    media_id INTEGER,
    episode_number INTEGER,
    tmdb_id INTEGER,

    FOREIGN KEY(season_id) REFERENCES tv_season(id),
    FOREIGN KEY(media_id) REFERENCES files(id)
);

CREATE TABLE IF NOT EXISTS media_subtitle (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    media_type TEXT CHECK(media_type IN ('movie','tv')) NOT NULL,
    media_id INTEGER NOT NULL,
    subtitle_id INTEGER NOT NULL,
    language TEXT,

    FOREIGN KEY(subtitle_id) REFERENCES files(id),
    UNIQUE(media_type, media_id, subtitle_id)
);

CREATE TABLE IF NOT EXISTS media_thumbnail (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    icon_id INTEGER,
    desktop_ini_id INTEGER,

    FOREIGN KEY(icon_id) REFERENCES files(id),
    FOREIGN KEY(desktop_ini_id) REFERENCES files(id)
);

CREATE TABLE IF NOT EXISTS media_genres (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    media_type TEXT CHECK(media_type IN ('movie','tv')),
    media_id INTEGER,
    genre_id INTEGER,
    FOREIGN KEY(genre_id) REFERENCES tmdb_genres(id)
);

CREATE TABLE IF NOT EXISTS tmdb_movie (
    id INTEGER PRIMARY KEY,
    title TEXT,
    original_title TEXT,
    original_language TEXT,
    overview TEXT,
    poster_path TEXT,
    release_date TEXT,
    genre_ids TEXT,
    vote_avg REAL,
    vote_count INTEGER,
    popularity REAL
);

CREATE TABLE IF NOT EXISTS tmdb_tv (
    id INTEGER PRIMARY KEY,
    title TEXT,
    original_title TEXT,
    original_language TEXT,
    overview TEXT,
    poster_path TEXT,
    first_air_date TEXT,
    genre_ids TEXT,
    vote_avg REAL,
    vote_count INTEGER,
    popularity REAL
);

CREATE TABLE IF NOT EXISTS tmdb_cast (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    media_type TEXT CHECK(media_type IN ('movie','tv')),
    media_id INTEGER,
    person_id INTEGER,
    character_name TEXT,
    cast_order INTEGER,
    FOREIGN KEY(person_id) REFERENCES tmdb_people(id)
);

CREATE TABLE IF NOT EXISTS tmdb_director (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    media_type TEXT CHECK(media_type IN ('movie','tv')),
    media_id INTEGER,
    person_id INTEGER,
    FOREIGN KEY(person_id) REFERENCES tmdb_people(id)
);

CREATE TABLE IF NOT EXISTS tmdb_people (
    id INTEGER PRIMARY KEY,
    name TEXT
);

CREATE TABLE IF NOT EXISTS tmdb_genres (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    media_type TEXT CHECK(media_type IN ('movie','tv')),
    name TEXT
);

CREATE TABLE IF NOT EXISTS tmdb_languages (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    iso_639_1 TEXT,
    english_name TEXT,
    name TEXT
);

