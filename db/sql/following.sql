-- name: FollowUser :exec
INSERT INTO following (following_id, follower_id) 
VALUES (@their_id, @my_id);


-- name: TestFollowUser :exec
INSERT INTO following (following_id, follower_id) 
VALUES ($1, $2)
ON CONFLICT DO NOTHING;

-- name: IsUserFollowingUser :one
SELECT COUNT(*)
FROM following
WHERE follower_id = $1 AND following_id = $2;


-- name: UnfollowUser :exec
DELETE FROM following
WHERE follower_id = $1 AND following_id = $2;


-- name: FollowUsername :exec
INSERT INTO following (following_id, follower_id)
VALUES (
    (
    SELECT id
    FROM users
    WHERE username = @their_username
    ), @my_id);

-- name: UnFollowUsername :exec
WITH their_id AS (
    SELECT id
    FROM users
    WHERE username = @their_username
),
delete_notification AS (
    DELETE FROM follow_notifications
    WHERE new_follower_id = @my_id 
          AND user_id = (SELECT id FROM their_id)
)
DELETE FROM following 
WHERE  following_id = (SELECT id FROM their_id)
       AND follower_id = @my_id;

-- name: AddFollowNotification :exec
INSERT INTO follow_notifications (
    user_id,
    new_follower_id
)
VALUES (
    (SELECT 
        id
    FROM users
    WHERE username = @their_username
    ),
    @my_id
);

-- DELETE FROM follow_notifications 


-- name: MarkAllFollowNotificationSeen :exec
UPDATE follow_notifications
SET seen = TRUE
WHERE user_id = @my_id;


-- name: GetFollowNotifications :many
SELECT 
    users.*
FROM follow_notifications
JOIN users
    ON follow_notifications.new_follower_id = users.id
WHERE follow_notifications.user_id = @my_id;


-- name: GetFollowing :many
SELECT 
    users.*
FROM following
JOIN users
ON following.following_id = users.id
WHERE  following.follower_id = (
	SELECT id
	FROM users
	WHERE users.username = @follower_username
);


-- name: GetFollowers :many
SELECT 
    users.*
FROM following
JOIN users
ON following.follower_id = users.id
WHERE following.following_id = (
	SELECT id
	FROM users
	WHERE users.username = @following_username
);


