-- name: AddStem :exec
INSERT INTO audio_file_stems (audio_file_id, stem_file_name)
VALUES ($1, $2);

-- name: GetStems :many
SELECT stem_file_name
FROM audio_file_stems
WHERE audio_file_id = $1;

-- name: GetNumberOfStemFiles :one
SELECT COUNT(*)
FROM audio_files
JOIN audio_file_stems
    ON audio_files.id = audio_file_stems.audio_file_id
WHERE audio_files.id = $1;
