-- name: SearchAudioFiles :many
WITH sc AS (
	SELECT 
		COUNT(*) AS stem_count,
		audio_file_id
	FROM audio_file_stems
	FULL JOIN audio_files
		ON audio_file_stems.audio_file_id = audio_files.id
	GROUP BY audio_file_stems.audio_file_id
)
SELECT	
	sqlc.embed(playables),
	sc.stem_count,
	audio_files.vis_arr,
	audio_files.musical_key,
	audio_files.musical_key_signature,
        audio_files.major_minor,
	stem_arrays.stem_file_names,
	ta.tag_array
FROM audio_files
JOIN users 
	ON audio_files.user_id = users.id
FULL JOIN audio_file_tag_arrays ta
	ON audio_files.id = ta.id
FULL JOIN sc
	ON audio_files.id = sc.audio_file_id
FULL JOIN stem_arrays
	ON audio_files.id = stem_arrays.af_id
JOIN playables
	ON audio_files.id = playables.id
JOIN search_doc
    ON search_doc.audio_file_id = audio_files.id
WHERE 
(
	(@bpm_radio::text = '') OR 
(@bpm_radio::text= 'use-exact' AND @exact_bpm::integer = -1 OR audio_files.bpm = @exact_bpm::integer) 
	OR (@bpm_radio::text = 'use-range' AND (@min_bpm::integer = -1 OR audio_files.bpm > @min_bpm::integer)
		   AND (@max_bpm::integer = -1 OR audio_files.bpm < @max_bpm::integer))
)
AND (
cast_to_musical_key(@musical_key::text) IS NULL OR audio_files.musical_key = cast_to_musical_key(@musical_key::text)
)
AND (
  cast_to_musical_key_signature(@key_sig::text) IS NULL OR audio_files.musical_key_signature = cast_to_musical_key_signature(@key_sig::text)
)
AND (
  cast_to_major_minor(@major_minor::text) IS NULL OR audio_files.major_minor = cast_to_major_minor(@major_minor::text)
)
AND (
  @stems_only::bool != 'on' OR sc.stem_count > 0
)
ORDER BY similarity(search_doc.doc, @keyword::text) DESC
LIMIT 20
OFFSET @page_offset::integer;


-- name: SearchForUsers :many 
SELECT 
   *
FROM user_search_rows
ORDER BY (similarity(user_search_rows.username, $1) ^ 2) + (similarity((COALESCE(user_search_rows.bio, '')), $1)) DESC
LIMIT 4;

-- name: TestingTagArrQuery :many
SELECT  audio_files.id AS af_id, array_agg(t.tag_name) AS tag_array
FROM audio_file_tags at
JOIN tags t
	ON t.id = at.tag_id
JOIN audio_files
ON audio_files.id = at.audio_file_id
GROUP BY audio_files.id;

-- name: SearchKeywordForDropdown :many
WITH ta AS (
	SELECT  audio_files.id AS af_id, array_agg(t.tag_name) AS tag_array
	FROM audio_file_tags at
	JOIN tags t
		ON t.id = at.tag_id
	JOIN audio_files
		ON audio_files.id = at.audio_file_id
	GROUP BY audio_files.id
)
SELECT 
	audio_files.title,
	users.username,
	audio_files.user_id,
	audio_files.artwork_src
FROM audio_files
JOIN users 
	ON audio_files.user_id = users.id
FULL JOIN ta
	ON audio_files.id = ta.af_id
JOIN search_doc
    ON search_doc.audio_file_id = audio_files.id
ORDER BY similarity(search_doc.doc, $1) DESC
LIMIT 5;



-- name: SearchKeywordWithFilters :many 
SELECT *
FROM audio_file_search_rows
WHERE 
(
	($3 = -1) OR 
	($3 = 1 AND $4 = -1 OR audio_files.bpm = $4) 
	OR ($3 = 2 AND ($5 = -1 OR audio_files.bpm > $5)
		   AND ($6 = -1 OR audio_files.bpm < $6))
)
AND (
  cast_to_musical_key($7) IS NULL OR audio_files.musical_key = cast_to_musical_key($7)
)
AND (
  cast_to_musical_key_signature($8) IS NULL OR audio_files.musical_key_signature = cast_to_musical_key_signature($8)
)
AND (
  cast_to_major_minor($9) IS NULL OR audio_files.major_minor = cast_to_major_minor($9)
)
AND (
  $10 != 'on' OR sc.stem_count > 0
)
ORDER BY similarity(search_doc.doc, $1) DESC
LIMIT 20
OFFSET $2;

-- $1 keyword
-- $2 page number
-- $3 bpmRadio
-- $4 exactBpm
-- $5 minBpm
-- $6 maxBpm
-- $7 musical key
-- $8 musical key signature
-- $9 major or minor
-- $10 includes stems ('on' or '')



