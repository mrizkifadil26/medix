CREATE INDEX IF NOT EXISTS idx_files_parent ON files(parent_type, parent_id);
CREATE INDEX IF NOT EXISTS idx_files_fingerprint ON files(fingerprint);

CREATE INDEX IF NOT EXISTS idx_movie_media ON movie(media_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_movie_unique
ON movie (title, year, source, "group", "collection");


CREATE INDEX IF NOT EXISTS idx_tv_season_tv ON tv_season(tv_id);
CREATE INDEX IF NOT EXISTS idx_tv_episode_season ON tv_episode(season_id);
CREATE INDEX IF NOT EXISTS idx_tv_episode_media ON tv_episode(media_id);

CREATE INDEX IF NOT EXISTS idx_subtitle_media ON media_subtitle(media_type, media_id);
CREATE INDEX IF NOT EXISTS idx_subtitle_file ON media_subtitle(subtitle_id);

CREATE INDEX IF NOT EXISTS idx_thumbnail_icon ON media_thumbnail(icon_id);

CREATE INDEX IF NOT EXISTS idx_media_genres_media ON media_genres(media_type, media_id);
CREATE INDEX IF NOT EXISTS idx_media_genres_genre ON media_genres(genre_id);

CREATE INDEX IF NOT EXISTS idx_tmdb_cast_media ON tmdb_cast(media_type, media_id);

CREATE INDEX IF NOT EXISTS idx_tmdb_director_media ON tmdb_director(media_type, media_id);

CREATE INDEX IF NOT EXISTS idx_tmdb_genres_name ON tmdb_genres(name);

CREATE INDEX IF NOT EXISTS idx_tmdb_languages_iso ON tmdb_languages(iso_639_1);
