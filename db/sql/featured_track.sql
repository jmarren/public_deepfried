-- name: GetFeaturedTrack :one
SELECT
    playables.id,
    playables.user_id,
    playables.audio_src,
    playables.username,
    title,
    bpm,
    playback_seconds,
    playables.created,
    playables.artwork_src
FROM playables
JOIN featured_track 
    ON featured_track.audio_file_id = playables.id
JOIN users 
    ON playables.user_id = users.id
LIMIT 1;


-- name: UpdateFeaturedTrack :exec
INSERT INTO featured_track (audio_file_id)
VALUES ($1)
ON CONFLICT (id) DO UPDATE
SET audio_file_id = $1;

-- CREATE TABLE IF NOT EXISTS featured_track (
--     id INT UNIQUE default(1),
--     user_id INT NOT NULL REFERENCES users,
--     audio_file_id INT NOT NULL REFERENCES audio_files,
--     Constraint CHK_single_row CHECK (id = 1)
-- );
--
