-- name: InsertTag :one
INSERT INTO tags (tag_name)
VALUES ($1)
ON CONFLICT DO NOTHING
RETURNING id;

-- WITH res AS (
--     SELECT COALESCE(COUNT(*), 0) as total
--     FROM tags
--     WHERE tags.tag_name = $1
-- ) 
-- SELECT 
--     CASE WHEN res.total > 0
--         THEN (
--             SELECT tags.id 
--             FROM tags
--             WHERE tags.tag_name = $1
--         )
--         ELSE 0
--     END
-- FROM res;

-- name: TestAddAudioFileTag :exec
INSERT INTO audio_file_tags (tag_id, audio_file_id)
VALUES ($1, $2)
ON CONFLICT DO NOTHING;

-- name: AddAudioFileTag :exec
INSERT INTO audio_file_tags (tag_id, audio_file_id)
VALUES ($1, $2);

-- name: GetTagId :one
SELECT tags.id
FROM tags
WHERE tag_name = @tag_name;



-- get all audio files and their tags: ******
-- SELECT tags.tag_name,audio_files.title
-- FROM tags
-- JOIN audio_file_tags 
-- ON audio_file_tags.tag_id = tags.id
-- JOIN audio_files
-- ON audio_files.id = audio_file_tags.audio_file_id;

-- name: GetAudioFileTags :many
SELECT tags.tag_name
FROM tags
JOIN audio_file_tags 
ON audio_file_tags.tag_id = tags.id
JOIN audio_files
ON audio_files.id = audio_file_tags.audio_file_id
WHERE audio_file_id = $1;



-- name: GetAudioFileTagsWithTitleAndUsername :many
SELECT tags.tag_name
FROM tags
JOIN audio_file_tags
    ON audio_file_tags.tag_id = tags.id
JOIN audio_files
    ON audio_files.id = audio_file_tags.audio_file_id
JOIN users  
    ON users.id = audio_files.user_id
WHERE audio_files.title = $1
    AND users.username = $2;


-- name: GetTagsOrderedByCount :many
SELECT tags.tag_name 
FROM tags
JOIN audio_file_tags 
ON audio_file_tags.tag_id = tags.id
JOIN audio_files
ON audio_files.id = audio_file_tags.audio_file_id
GROUP BY tags.tag_name
ORDER BY COUNT(*) DESC
LIMIT 30;


-- name: GetTagCount :one
SELECT 
    COUNT(*)
FROM tags
WHERE tag_name = @tag_name;


--
-- SELECT tags.tag_name
-- FROM tags
-- JOIN audio_file_tags 
-- ON audio_file_tags.tag_id = tags.id
-- JOIN audio_files
-- ON audio_files.id = audio_file_tags.audio_file_id
-- WHERE audio_file_id = $1;
