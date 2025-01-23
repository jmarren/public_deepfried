

-- name: AddUserDownload :exec
INSERT INTO user_downloads (user_id, audio_file_id)
VALUES (@user_id, @audio_id);


-- name: GetUserDownloads :many
SELECT 
    playables.*
FROM playables
JOIN user_downloads
    ON user_downloads.audio_file_id = playables.id
WHERE user_downloads.user_id = $1
LIMIT 40;



-- name: GetUserDownloadsWithKeyword :many
SELECT 
    playables.*
FROM playables
JOIN user_downloads
    ON user_downloads.audio_file_id = playables.id
WHERE user_downloads.user_id = $1
ORDER BY ((similarity(playables.title, @keyword) ^ 2) + similarity(playables.username, @keyword)) DESC
LIMIT 40;

-- (similarity(user_search_rows.username, $1) ^ 2) + (similarity((COALESCE(user_search_rows.bio, '')), $1)) DESC

