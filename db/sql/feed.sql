-- name: GetUserFeed :many
SELECT 
    sqlc.embed(playables),
    audio_file_tag_arrays.tag_array,
    audio_files.vis_arr,
    audio_files.usage_rights
FROM users
JOIN following
	ON following.following_id = users.id
JOIN playables
	ON following.following_id = playables.user_id
FULL JOIN audio_file_tag_arrays 
	ON playables.id = audio_file_tag_arrays.id  
JOIN audio_files 
	ON playables.id = audio_files.id
WHERE following.follower_id = @my_id;





