-- name: InsertTestPin :exec
INSERT INTO pins (user_id, file_id)
VALUES ($1, $2)
ON CONFLICT DO NOTHING; -- noqa

-- name: GetUserPins :many
WITH audio_ids AS (
    SELECT file_id
    FROM pins
    WHERE pins.user_id = $1
)
SELECT 
        sqlc.embed(playables),
        audio_file_tag_arrays.tag_array
FROM playables
JOIN audio_file_tag_arrays
        ON playables.id = audio_file_tag_arrays.id
JOIN audio_ids
        ON playables.id = audio_ids.file_id
LIMIT 4;
-- WITH audio_ids AS (
--     SELECT file_id
--     FROM pins
--     WHERE pins.user_id = $1
-- )
-- SELECT 
--     id,
--     title,
--     artwork_src,
--     audio_src
-- FROM audio_files
-- INNER JOIN audio_ids ON audio_ids.file_id = audio_files.id
-- LIMIT 4;


-- name: DeleteAllUserPins :exec
DELETE FROM pins
WHERE pins.user_id = $1;

-- name: EditUserPins :exec
INSERT INTO pins (user_id, file_id) 
SELECT 
        audio_files.user_id,
        audio_files.id
FROM 
        audio_files
WHERE audio_files.title = $2
    AND audio_files.user_id = $1;


/*
WITH ta AS (
	SELECT  audio_files.id AS af_id, array_agg(t.tag_name) AS tag_array
	FROM audio_file_tags at
	JOIN tags t
		ON t.id = at.tag_id
	JOIN audio_files
		ON audio_files.id = at.audio_file_id
	GROUP BY audio_files.id
), audio_ids AS (
    SELECT file_id
    FROM pins
    WHERE pins.user_id = $1
)
SELECT 
    audio_files.id,
    title,
    artwork_src,
    audio_src,
    ta.tag_array
FROM audio_files
INNER JOIN audio_ids ON audio_ids.file_id = audio_files.id
FULL JOIN ta
	ON audio_files.id = ta.af_id
LIMIT 4;
*/

