-- name: AddAudioFile :one
INSERT INTO audio_files(
    user_id,
    title,
    audio_src,
    bpm,
    musical_key,
    musical_key_signature,
    major_minor,
    playback_seconds,
    file_size,
    usage_rights,
    artwork_src,
    vis_arr
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
RETURNING id;

-- name: AddTestAudioFile :one
WITH res AS (
    INSERT INTO audio_files (
        user_id,
        title,
        audio_src,
        bpm,
        musical_key,
        musical_key_signature,
        major_minor,
        usage_rights,
        playback_seconds,
        file_size,
        vis_arr,
        artwork_src
    )
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
    ON CONFLICT DO NOTHING
    RETURNING id --noqa
)
SELECT id FROM res
    UNION ALL
SELECT id FROM audio_files WHERE user_id = $1 AND title = $2
LIMIT 1;

-- name: GetUserAudioFiles :many
SELECT
    playables.*
FROM playables
JOIN audio_files 
    ON audio_files.id = playables.id
WHERE playables.user_id = $1
ORDER BY audio_files.listen_count DESC
LIMIT 20;

-- name: GetUserAudioFilesWithLimit :many
SELECT
    sqlc.embed(playables),
    audio_file_tag_arrays.tag_array
FROM playables
FULL JOIN audio_file_tag_arrays
        ON playables.id = audio_file_tag_arrays.id
JOIN audio_files
    ON playables.id = audio_files.id 
WHERE playables.user_id = @user_id
ORDER BY audio_files.listen_count DESC
LIMIT @number_of_results;


-- name: GetFourUserAudioFiles :many
SELECT
    id,
    title,
    audio_src,
    bpm,
    playback_seconds,
    created,
    vis_arr,
    artwork_src
FROM audio_files
WHERE user_id = $1
LIMIT 4;

-- name: GetMostAudioFilesOrderedByListen :many
SELECT
    a.id,
    a.title,
    a.audio_src,
    a.bpm,
    a.playback_seconds,
    a.created,
    a.vis_arr,
    a.user_id,
    u.username,
    a.artwork_src
FROM audio_files AS a
INNER JOIN users AS u
    ON a.user_id = u.id
ORDER BY a.listen_count DESC
LIMIT 20;


-- name: DoesTitleExistForUser :one
SELECT COUNT(audio_files.title)
FROM audio_files
JOIN users
    ON users.id = audio_files.user_id
WHERE 
    audio_files.title = $1 
    AND users.id = $2;

-- name: GetMostPopularAudioFiles :many
SELECT
    playables.*
FROM playables
JOIN audio_files 
    ON audio_files.id = playables.id
ORDER BY audio_files.listen_count DESC
LIMIT 20;


-- name: GetEditorsPicks :many
SELECT
    playables.*
FROM playables
JOIN editors_picks 
    ON playables.id = editors_picks.audio_file_id
LIMIT 20;


-- name: GetPlayableByTitleAndUsername :one
SELECT 
    sqlc.embed(playables),
    audio_files.vis_arr
FROM 
    playables
JOIN audio_files
    ON playables.id = audio_files.id
WHERE 
    playables.username = $1
    AND playables.title = $2;


-- name: GetAudioFile :one 
SELECT 
    a.id,
    a.title,
    a.audio_src,
    a.bpm,
    a.playback_seconds,
    a.created,
    a.vis_arr,
    a.usage_rights,
    a.user_id,
    u.username,
    a.artwork_src,
    a.musical_key,
    a.major_minor,
    a.musical_key_signature
FROM audio_files AS a
INNER JOIN users AS u
    ON a.user_id = u.id
WHERE u.username = $1
    AND a.title = $2
LIMIT 1;


-- name: ListUserAudioFilesAndWhetherPinned :many
WITH uaf AS (
    SELECT 
        audio_files.id as fid,
        audio_files.title as ft
    FROM audio_files
    WHERE audio_files.user_id = $1
)
SELECT 
    uaf.ft AS title, 
    CASE 
    WHEN (
        SELECT COUNT(*) 
        FROM pins 
        WHERE pins.file_id = uaf.fid
    ) > 0 THEN TRUE
    ELSE FALSE
    END AS pinned
FROM uaf;







-- name: GetJustAdded :many
SELECT
    playables.*
FROM playables
ORDER BY playables.created DESC
LIMIT 4;

-- name: GetAudioFileById :one
SELECT 
    playables.*
FROM playables
WHERE playables.id = @audio_file_id;


