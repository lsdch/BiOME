CREATE TYPE vocab_resolution AS (
    id uuid,
    code text,
    label text,
    confidence float,
    match_type text,
    should_accept boolean
);

CREATE OR REPLACE FUNCTION vocab_decision(
        v_id uuid,
        v_code text,
        v_label text,
        v_score float,
        min_auto float
    ) RETURNS vocab_resolution LANGUAGE sql AS $$
SELECT v_id,
    v_code,
    v_label,
    COALESCE(v_score, 0),
    CASE
        WHEN v_score = 1 THEN 'exact'
        WHEN v_score >= min_auto THEN 'fuzzy'
        ELSE 'none'
    END,
    v_score >= min_auto;
$$;