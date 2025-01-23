DROP SCHEMA public CASCADE;
CREATE SCHEMA public;

CREATE TYPE MUSICAL_KEY AS ENUM ('C', 'D', 'E', 'F', 'G', 'A', 'B');
CREATE TYPE MUSICAL_KEY_SIGNATURE AS ENUM ('flat', 'natural', 'sharp');
CREATE TYPE MAJOR_MINOR AS ENUM ('Major', 'Minor');

CREATE TABLE IF NOT EXISTS musical_key_lookup (
    enum_val MUSICAL_KEY PRIMARY KEY NOT NULL,
    text_val TEXT NOT NULL
);

INSERT INTO musical_key_lookup (enum_val, text_val) VALUES ('C', 'C')
ON CONFLICT DO NOTHING;

INSERT INTO musical_key_lookup (enum_val, text_val) VALUES ('D', 'D')
ON CONFLICT DO NOTHING;


INSERT INTO musical_key_lookup (enum_val, text_val) VALUES ('E', 'E')
ON CONFLICT DO NOTHING;


INSERT INTO musical_key_lookup (enum_val, text_val) VALUES ('F', 'F')
ON CONFLICT DO NOTHING;


INSERT INTO musical_key_lookup (enum_val, text_val) VALUES ('G', 'G')
ON CONFLICT DO NOTHING;


INSERT INTO musical_key_lookup (enum_val, text_val) VALUES ('A', 'A')
ON CONFLICT DO NOTHING;


INSERT INTO musical_key_lookup (enum_val, text_val) VALUES ('B', 'B')
ON CONFLICT DO NOTHING;
