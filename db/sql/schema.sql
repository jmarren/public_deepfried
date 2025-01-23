-- GRANT CREATE ON SCHEMA public TO PUBLIC;
CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE EXTENSION IF NOT EXISTS PLPGSQL;
-- CREATE EXTENSION IF NOT EXISTS "uuid-oosp";


-- name: InitUsers :exec
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    cognito_id VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(20) UNIQUE NOT NULL
);


CREATE TABLE IF NOT EXISTS profiles (
    user_id UUID PRIMARY KEY NOT NULL 
        REFERENCES users(id)
        ON DELETE CASCADE,
    created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    user_location VARCHAR(50),
    last_modified TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    bio VARCHAR(100)
);

-- name: InitAudioFiles :exec
CREATE TABLE IF NOT EXISTS audio_files ( -- noqa
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    user_id UUID NOT NULL 
        REFERENCES users
        ON DELETE CASCADE,
    title VARCHAR(20) NOT NULL,
    audio_src VARCHAR(255) NOT NULL,
    bpm INT NOT NULL,
    musical_key MUSICAL_KEY NOT NULL,
    musical_key_signature MUSICAL_KEY_SIGNATURE NOT NULL,
    major_minor MAJOR_MINOR NOT NULL,
    playback_seconds INT NOT NULL,
    file_size INTEGER NOT NULL,
    created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    listen_count INTEGER NOT NULL DEFAULT 0,
    download_count INTEGER NOT NULL DEFAULT 0,
    vis_arr INTEGER [] NOT NULL,
    usage_rights VARCHAR(255),
    artwork_src VARCHAR(255) NOT NULL,
    last_modified TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (user_id, title)
);


CREATE TABLE IF NOT EXISTS tags (
    id SERIAL PRIMARY KEY NOT NULL,
    tag_name VARCHAR(25) NOT NULL,
    UNIQUE (tag_name)
);

CREATE TABLE IF NOT EXISTS audio_file_tags (
    tag_id INT NOT NULL 
        REFERENCES tags
        ON DELETE CASCADE,
    audio_file_id UUID NOT NULL 
        REFERENCES audio_files
        ON DELETE CASCADE,
    PRIMARY KEY (tag_id, audio_file_id)
);

CREATE TABLE IF NOT EXISTS audio_file_stems (
    audio_file_id UUID NOT NULL
        REFERENCES audio_files
        ON DELETE CASCADE,
    stem_file_name VARCHAR(255) NOT NULL,
    PRIMARY KEY (audio_file_id, stem_file_name)
);

-- name: InitPins :exec
CREATE TABLE IF NOT EXISTS pins (
    id SERIAL PRIMARY KEY,
    file_id UUID NOT NULL 
        REFERENCES audio_files
        ON DELETE CASCADE,
    user_id UUID NOT NULL 
        REFERENCES users
        ON DELETE CASCADE,
    UNIQUE (user_id, file_id)
);

-- name: InitFeaturedSection :exec
CREATE TABLE IF NOT EXISTS featured_track (
    id INT UNIQUE default(1),
    audio_file_id UUID NOT NULL 
        REFERENCES audio_files
        ON DELETE CASCADE,
    Constraint CHK_single_row CHECK (id = 1)
);


-- CREATE TABLE IF NOT EXISTS editors_picks (
--     id INT UNIQUE default(1),
--     audio_file_id UUID NOT NULL 
--         REFERENCES audio_files
--         ON DELETE CASCADE,
--     Constraint CHK_fifteen_rows_or_less CHECK (id < 16)
-- );


-- name: InitFollowing :exec
CREATE TABLE IF NOT EXISTS following (
    follower_id UUID NOT NULL 
        REFERENCES users
        ON DELETE CASCADE,
    following_id UUID NOT NULL 
    REFERENCES users
    ON DELETE CASCADE,
    PRIMARY KEY (follower_id, following_id)
);


CREATE TABLE IF NOT EXISTS user_downloads (
    user_id UUID NOT NULL,
    audio_file_id UUID NOT NULL,
    PRIMARY KEY (user_id, audio_file_id)
);

CREATE TABLE IF NOT EXISTS search_doc (
    doc TEXT,
    user_id UUID NOT NULL 
        REFERENCES users
        ON DELETE CASCADE,
    audio_file_id UUID NOT NULL 
        REFERENCES audio_files
        ON DELETE CASCADE, 
    PRIMARY KEY(user_id, audio_file_id)
);

CREATE TABLE IF NOT EXISTS follow_notifications (
    user_id UUID NOT NULL
        REFERENCES users,
    new_follower_id UUID NOT NULL
        REFERENCES users,
    seen BOOL DEFAULT FALSE,
    PRIMARY KEY (user_id, new_follower_id)
);


CREATE TABLE IF NOT EXISTS submission_postings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    creator_id UUID NOT NULL
        REFERENCES users,
    title varchar(255) NOT NULL,
    description TEXT
);

CREATE TABLE IF NOT EXISTS submission_entries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    user_id UUID NOT NULL REFERENCES users,
    posting_id UUID NOT NULL REFERENCES submission_postings,
    audio_file_id UUID NOT NULL REFERENCES audio_files
);



CREATE TABLE IF NOT EXISTS admins (
    user_id UUID PRIMARY KEY NOT NULL
        REFERENCES users
);


CREATE TABLE IF NOT EXISTS editors_picks (
    id SERIAL PRIMARY KEY NOT NULL,
    audio_file_id UUID NOT NULL 
        REFERENCES audio_files
);


CREATE OR REPLACE FUNCTION fn_concat_tags(a UUID) 
    RETURNS text 
    LANGUAGE PLPGSQL
AS
$$
        DECLARE 
        arow RECORD;
        tags_concat TEXT := '';
        BEGIN
        FOR arow IN (
            SELECT tag_name AS t
            FROM tags
            JOIN audio_file_tags
                ON audio_file_tags.tag_id = tags.id
            WHERE audio_file_tags.audio_file_id = a
        )
        LOOP
            tags_concat := CONCAT(tags_concat , ' ' , arow.t);
            RAISE NOTICE 'tags_concat %', tags_concat;
        END LOOP;
        RETURN tags_concat;
    END;
$$;

CREATE OR REPLACE FUNCTION fn_update_search_doc_new_audio_file()
    RETURNS TRIGGER
    LANGUAGE PLPGSQL
AS
$$
    DECLARE 
        doc TEXT := '';
    BEGIN
    INSERT INTO search_doc (doc, user_id, audio_file_id)
        VALUES (
            (
                SELECT CONCAT(COALESCE(fn_concat_tags(NEW.id), '') , ' ' , users.username, ' ', NEW.title)
                FROM users
                WHERE NEW.user_id = users.id
            ),
            NEW.user_id,
            NEW.id
            );
        RETURN NEW;
        END;
$$;


CREATE OR REPLACE TRIGGER new_audio_file AFTER INSERT OR UPDATE ON audio_files
    FOR EACH ROW EXECUTE FUNCTION fn_update_search_doc_new_audio_file();

CREATE OR REPLACE FUNCTION fn_update_search_doc_new_audio_tag()
    RETURNS TRIGGER
    LANGUAGE PLPGSQL
AS
$$
    DECLARE 
        doc TEXT := '';
    BEGIN
    INSERT INTO search_doc (doc, user_id, audio_file_id)
        VALUES (
            (
                SELECT CONCAT(fn_concat_tags(NEW.audio_file_id) , ' ' , users.username, ' ', audio_files.title)
                FROM audio_files
                JOIN users ON users.id = audio_files.user_id
                WHERE audio_files.id = NEW.audio_file_id
            ),
            (
                SELECT users.id
                FROM audio_files
                JOIN audio_file_tags 
                        ON audio_files.id = audio_file_tags.audio_file_id
                JOIN users 
                        ON audio_files.user_id = users.id
                WHERE audio_files.id = NEW.audio_file_id
                LIMIT 1
            ),
            NEW.audio_file_id
            )
        ON CONFLICT (user_id, audio_file_id)
        DO UPDATE SET doc = EXCLUDED.doc;
        RETURN NEW;
        END;
$$;

CREATE OR REPLACE TRIGGER new_audio_tag AFTER INSERT OR UPDATE ON audio_file_tags
    FOR EACH ROW EXECUTE FUNCTION fn_update_search_doc_new_audio_tag();


CREATE OR REPLACE FUNCTION cast_to_musical_key(t_input text) 
RETURNS MUSICAL_KEY AS $$
DECLARE t_enum_val MUSICAL_KEY DEFAULT NULL;
BEGIN
    BEGIN
        t_enum_val := t_input::MUSICAL_KEY;
    EXCEPTION WHEN OTHERS THEN 
    RETURN NULL;
    END;
RETURN t_enum_val;
END;
$$ LANGUAGE PLPGSQL IMMUTABLE;


CREATE OR REPLACE FUNCTION cast_to_musical_key_signature(t_input text) 
RETURNS MUSICAL_KEY_SIGNATURE AS $$
DECLARE t_enum_val MUSICAL_KEY_SIGNATURE DEFAULT NULL;
BEGIN
    BEGIN
        t_enum_val := t_input::MUSICAL_KEY_SIGNATURE;
    EXCEPTION WHEN OTHERS THEN 
    RETURN NULL;
    END;
RETURN t_enum_val;
END;
$$ LANGUAGE PLPGSQL IMMUTABLE;


CREATE OR REPLACE FUNCTION cast_to_major_minor(t_input text) 
RETURNS MAJOR_MINOR AS $$
DECLARE t_enum_val MAJOR_MINOR DEFAULT NULL;
BEGIN
    BEGIN
        t_enum_val := t_input::MAJOR_MINOR;
    EXCEPTION WHEN OTHERS THEN 
    RETURN NULL;
    END;
RETURN t_enum_val;
END;
$$ LANGUAGE PLPGSQL IMMUTABLE;

CREATE OR REPLACE VIEW basic_users AS (
SELECT id, username
FROM users
);



CREATE OR REPLACE VIEW carousel_cards AS (
    SELECT
        users.username,
        user_id,
        title,
        audio_src,
        bpm,
        playback_seconds,
        audio_files.created,
        artwork_src,
        listen_count
    FROM audio_files
    JOIN users
        ON audio_files.user_id = users.id
);






CREATE OR REPLACE VIEW audio_file_search_rows AS (
WITH ta AS (
	SELECT  audio_files.id AS af_id, array_agg(t.tag_name) AS tag_array
	FROM audio_file_tags at
	JOIN tags t
		ON t.id = at.tag_id
	JOIN audio_files
		ON audio_files.id = at.audio_file_id
	GROUP BY audio_files.id
), sc AS (
	SELECT 
		COUNT(*) AS stem_count,
		audio_file_id
	FROM audio_file_stems
	FULL JOIN audio_files
		ON audio_file_stems.audio_file_id = audio_files.id
	GROUP BY audio_file_stems.audio_file_id
),
sa  AS (
	SELECT audio_files.id AS af_id, array_agg(stem_file_name) AS stem_file_names
	FROM audio_file_stems
	JOIN audio_files
		ON audio_file_stems.audio_file_id = audio_files.id
	GROUP BY audio_files.id
)
SELECT 
	audio_files.title,
	users.username,
	sc.stem_count,
	audio_files.audio_src,
	audio_files.user_id,
	audio_files.artwork_src,
	audio_files.bpm,
	audio_files.vis_arr,
	audio_files.playback_seconds,
	audio_files.created,
	audio_files.musical_key,
	audio_files.musical_key_signature,
        audio_files.major_minor,
	sa.stem_file_names,
	ta.tag_array,
        search_doc.doc
FROM audio_files
JOIN users 
	ON audio_files.user_id = users.id
FULL JOIN ta
	ON audio_files.id = ta.af_id
FULL JOIN sc
	ON audio_files.id = sc.audio_file_id
FULL JOIN sa 
	ON audio_files.id = sa.af_id
JOIN search_doc
    ON search_doc.audio_file_id = audio_files.id
);


CREATE OR REPLACE VIEW audio_and_artwork AS (
    SELECT id,
           artwork_src,
           audio_src
    FROM audio_files
);


CREATE OR REPLACE VIEW follower_counts AS (
    SELECT 
            following_id,
            COUNT(*)
    FROM following
    GROUP BY (following_id)
);

CREATE OR REPLACE VIEW following_counts AS (
    SELECT 
            follower_id,
            COUNT(*)
    FROM following
    GROUP BY (follower_id)
);


CREATE OR REPLACE VIEW audio_file_tag_arrays AS (
    SELECT  
            audio_files.id,
            array_agg(t.tag_name) AS tag_array
    FROM audio_file_tags at
    JOIN tags t
            ON t.id = at.tag_id
    JOIN audio_files
            ON audio_files.id = at.audio_file_id
    GROUP BY audio_files.id
);


CREATE OR REPLACE VIEW user_pins AS (
SELECT 
	pins.user_id,
	pins.file_id AS audio_file_id,
	audio_file_tag_arrays.tag_array,
        audio_files.title,
	audio_files.audio_src,
	audio_files.artwork_src
FROM pins
JOIN audio_files 
	ON audio_files.id = pins.file_id
FULL JOIN audio_file_tag_arrays
	ON audio_file_tag_arrays.id = pins.file_id
WHERE pins.user_id IS NOT NULL
);

CREATE OR REPLACE VIEW user_search_rows AS (
SELECT 
    username,
    bio,
    id
FROM users
JOIN profiles
    ON users.id = profiles.user_id
);


CREATE OR REPLACE VIEW play_sources AS (
SELECT 
    audio_files.id,
    audio_files.audio_src,
    audio_files.artwork_src
FROM audio_files
);


CREATE OR REPLACE VIEW playables AS (
SELECT 
    audio_files.id,
    audio_files.user_id,
    audio_files.audio_src,
    users.username,
    audio_files.title,
    audio_files.bpm,
    audio_files.playback_seconds,
    audio_files.created,
    audio_files.artwork_src
FROM audio_files
JOIN users
ON users.id = audio_files.user_id
);


CREATE OR REPLACE VIEW user_following AS (
    SELECT 
        following.following_id AS fi,
        users.username AS fu
    FROM 
        following
    JOIN users
        ON following.following_id = users.id
);



CREATE OR REPLACE VIEW stem_arrays AS (
	SELECT audio_files.id AS af_id, array_agg(stem_file_name) AS stem_file_names
	FROM audio_file_stems
	JOIN audio_files
		ON audio_file_stems.audio_file_id = audio_files.id
	GROUP BY audio_files.id
);
