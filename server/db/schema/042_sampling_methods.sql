CREATE TABLE IF NOT EXISTS sampling_methods (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    code CITEXT NOT NULL UNIQUE,
    name CITEXT NOT NULL UNIQUE,
    description TEXT
);


CREATE TABLE IF NOT EXISTS events_sampling_methods (
    sampling_id UUID NOT NULL REFERENCES samplings (id) ON DELETE CASCADE,
    method_id UUID NOT NULL REFERENCES sampling_methods (id) ON DELETE CASCADE,
    PRIMARY KEY (sampling_id, method_id)
);

CREATE OR REPLACE FUNCTION resolve_method(
        input_text text,
        min_auto float DEFAULT 0.60
    ) RETURNS SETOF vocab_resolution LANGUAGE plpgsql AS $$
DECLARE v_clean text;
DECLARE v_result vocab_resolution;
BEGIN v_clean := lower(trim(input_text));

    ------------------------------------------------------------------
-- EXACT MATCH
------------------------------------------------------------------
SELECT id,
    code,
    name,
    1.0 INTO v_result
FROM sampling_methods
WHERE lower(code) = v_clean
    OR lower(name) = v_clean
LIMIT 1;

    IF FOUND THEN RETURN NEXT vocab_decision(
    v_result.id,
    v_result.code,
    v_result.label,
    1.0,
    min_auto
);
RETURN;
END IF;

------------------------------------------------------------------
-- FUZZY MATCH
------------------------------------------------------------------
SELECT m.id,
    m.code,
    m.name,
    similarity(m.name, v_clean) INTO v_result
FROM sampling_methods m
WHERE m.name % v_clean
ORDER BY similarity(m.name, v_clean) DESC
LIMIT 1;

    IF FOUND THEN RETURN NEXT vocab_decision(
    v_result.id,
    v_result.code,
    v_result.label,
    v_result.confidence,
    min_auto
);
RETURN;
END IF;

    ------------------------------------------------------------------
-- FALLBACK RAW
------------------------------------------------------------------
RETURN NEXT (
    NULL::uuid,
    NULL::text,
    input_text,
    0.0,
    'none',
    false
);

END;
$$;