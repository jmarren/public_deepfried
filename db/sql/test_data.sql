
INSERT INTO users (
    cognito_id,
    username
)
VALUES (
    '00001',
    'johnmarren'
)
ON CONFLICT DO NOTHING; -- noqa

INSERT INTO users (
    cognito_id,
    username
)
VALUES (
    '00002',
    'wonderlust'
)
ON CONFLICT DO NOTHING;

INSERT INTO users (
    cognito_id,
    username
)
VALUES (
    '00003',
    'lovejames'
)
ON CONFLICT DO NOTHING;

INSERT INTO users (
    cognito_id,
    username
)
VALUES (
    '00004',
    'blakefoster'
)
ON CONFLICT DO NOTHING;

INSERT INTO audio_files (
    user_id,
    title,
    playback_time,
    file_size,
    artwork,
    vis_arr
)
VALUES (
    (SELECT id FROM users WHERE username = 'wonderlust' LIMIT 1),
    'the more the merrier',
    '1:01:32',
    123456,
    'https://loop.com/artwork/1/the-more-the-merrier',
    '{1, 2, 3, 4, 44}'
)
ON CONFLICT DO NOTHING;

INSERT INTO audio_files (
    user_id,
    title,
    playback_time,
    file_size,
    artwork,
    vis_arr
)
VALUES (
    (SELECT id FROM users WHERE username = 'lovejames'),
    'quantum leap',
    '3:08:32',
    4321,
    'https://loop.com/artwork/2/quantum-leap',
    '{91, 20, 40, 83, 21}'
)
ON CONFLICT DO NOTHING;

INSERT INTO audio_files (
    user_id,
    title,
    playback_time,
    file_size,
    artwork,
    vis_arr
)
VALUES (
    (SELECT id FROM users WHERE username = 'blakefoster'),
    'painstaking',
    '2:38:53',
    4321,
    'https://loop.com/artwork/2/painstaking',
    '{48, 0, 4, 83, 88, 64, 73, 82, 28, 74, 73, 82, 73, 18, 19, 54, 58}'
)
ON CONFLICT DO NOTHING;


INSERT INTO audio_files (
    user_id,
    title,
    playback_time,
    file_size,
    artwork,
    vis_arr
)
VALUES (
    (SELECT id FROM users WHERE username = 'blakefoster'),
    'rat race',
    '2:19:34',
    983013,
    'https://loop.com/artwork/2/rat_race',
    '{0, 0, 4, 1, 3, 64, 73, 82, 28, 83, 12, 11, 10, 18, 19, 33, 58}'
)
ON CONFLICT DO NOTHING;


