
-- name: CreateUser :one
INSERT INTO users (cognito_id, username)
VALUES (@cognito_id, @username)
RETURNING id as inserted_id;


-- name: CreateProfile :exec
INSERT INTO profiles (user_id, bio)
VALUES (@user_id::uuid, @bio);

-- name: CreateUserTest :one
INSERT INTO users (cognito_id, username)
VALUES (@cognito_id, @username)
ON CONFLICT DO NOTHING
RETURNING id as inserted_id;

-- name: CreateProfileTest :exec
INSERT INTO profiles (user_id, bio)
VALUES (@user_id, @bio)
ON CONFLICT DO NOTHING;



-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

-- name: GetUser :one
SELECT
    id,
    cognito_id,
    username
FROM users
WHERE id = $1;


-- name: IsCognitoIdPresent :one
SELECT
    CASE
        WHEN
            EXISTS (
                SELECT 1
                FROM users
                WHERE cognito_id = $1
            )
            THEN 1
        ELSE 0
    END;


-- name: GetUserWithCognitoId :one
SELECT
    *
FROM
    users
WHERE
    cognito_id = $1
LIMIT 1;


-- name: GetUserBio :one 
SELECT bio
FROM profiles
WHERE user_id = $1;


-- name: GetUserWithUsername :one
SELECT
    *
FROM 
    users
WHERE
    username = $1
LIMIT 1;


-- name: DoesUsernameExist :one 
SELECT COUNT(*)
FROM users
WHERE username = $1;

-- name: UpdateUserBio :exec
UPDATE profiles
SET bio = $2
WHERE user_id = $1;

-- name: UpdateUserUsername :exec
UPDATE users
SET username = $2
WHERE id = $1;


-- name: GetUserInfo :one 
WITH uid AS (
        SELECT id AS i
        FROM users
        WHERE username = $1
)
SELECT 
    (
        SELECT bio FROM profiles WHERE user_id = uid.i
    ),
    (
    SELECT 
    COUNT(*) AS followers
    FROM following
    WHERE following_id = uid.i
    ),
    (
    SELECT 
    COUNT(*) AS following
    FROM following
    WHERE follower_id = uid.i
    )
FROM users
JOIN uid 
    ON uid.i = users.id
WHERE users.username = $1;


-- name: GetNumberFollowing :one 
SELECT count
FROM following_counts
WHERE follower_id = $1;


-- name: GetNumberFollowers :one 
SELECT count
FROM follower_counts
WHERE following_id = $1;


-- name: GetAmIFollowing :one
SELECT CASE WHEN COALESCE(COUNT(*), 0) > 0 THEN TRUE
        ELSE FALSE
    END 
FROM following
WHERE following.follower_id = @my_id
    AND following.following_id = @their_id;


-- name: IsUserAdmin :one
SELECT CASE WHEN COALESCE(COUNT(*), 0) > 0 THEN TRUE
            ELSE FALSE
        END
FROM admins 
WHERE admins.user_id = @user_id;


-- name: GetFiveNewestUsers :many
SELECT * 
FROM users
JOIN profiles
    ON users.id = profiles.user_id
ORDER BY profiles.created DESC
LIMIT 5;










