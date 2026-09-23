-- Habitat groups and habitat tags
CREATE TABLE habitat_groups (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
	label CITEXT NOT NULL,
	description TEXT,
	exclusive_elements BOOLEAN NOT NULL DEFAULT true,
	CONSTRAINT habitat_group_label_unique UNIQUE (label),
	CONSTRAINT habitat_group_label_not_empty CHECK (btrim(label) <> '')
);


-- Habitats belong to a group and may form a hierarchy (parent)
CREATE TABLE habitats (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
	label CITEXT NOT NULL,
	description TEXT,
	habitat_group_id UUID NOT NULL REFERENCES habitat_groups (id) ON DELETE CASCADE,
	CONSTRAINT habitat_label_not_empty CHECK (btrim(label) <> ''),
	CONSTRAINT habitat_description_length CHECK (char_length(coalesce(description, '')) <= 4000),
	CONSTRAINT uq_habitat_label UNIQUE (label)
);

ALTER TABLE habitat_groups
ADD COLUMN parent_habitat_id UUID REFERENCES habitats (id) ON DELETE
SET NULL;

CREATE INDEX habitat_group_idx ON habitats (habitat_group_id);

-- Full-text search document column (label + description)
-- ALTER TABLE habitats
-- ADD COLUMN document tsvector GENERATED ALWAYS AS (
-- 		to_tsvector(
-- 			'simple',
-- 			coalesce(label, '') || ' ' || coalesce(description, '')
-- 		)
-- 	) STORED;
-- CREATE INDEX habitat_document_idx ON habitats USING GIN (document);
-- Association table linking samplings to habitats
CREATE TABLE samplings_habitats (
	sampling_id UUID NOT NULL REFERENCES samplings (id) ON DELETE CASCADE,
	habitat_id UUID NOT NULL REFERENCES habitats (id) ON DELETE CASCADE,
	PRIMARY KEY (sampling_id, habitat_id)
);
CREATE INDEX samplings_habitats_habitat_idx ON samplings_habitats (habitat_id);


--
-- Prevent cycles in parent relationship
--
CREATE OR REPLACE FUNCTION validate_group_parent_not_in_subtree()
RETURNS trigger AS $$
DECLARE
    group_to_check UUID;
    is_invalid BOOLEAN;
BEGIN
    -- AFTER triggers see the new parent/group, including multi-row updates.
    IF TG_TABLE_NAME = 'habitat_groups' THEN
        group_to_check := NEW.id;
    ELSE
        group_to_check := NEW.habitat_group_id;
    END IF;

    -- A group's parent is a habitat, whose group is the next ancestor.
    -- UNION also terminates traversal if an existing cycle is encountered.
    WITH RECURSIVE ancestors AS (
        SELECT h.habitat_group_id AS id
        FROM habitat_groups g
            JOIN habitats h ON h.id = g.parent_habitat_id
        WHERE g.id = group_to_check
        UNION
        SELECT h.habitat_group_id
        FROM ancestors a
            JOIN habitat_groups g ON g.id = a.id
            JOIN habitats h ON h.id = g.parent_habitat_id
    )
    SELECT EXISTS (
        SELECT 1 FROM ancestors WHERE id = group_to_check
    ) INTO is_invalid;

    IF is_invalid THEN
        RAISE EXCEPTION 'Invalid habitat hierarchy: a group cannot depend on one of its own descendant habitats'
            USING ERRCODE = 'HB001';
    END IF;

    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER habitat_groups_validate_parent
AFTER INSERT OR UPDATE OF parent_habitat_id ON habitat_groups
FOR EACH ROW EXECUTE FUNCTION validate_group_parent_not_in_subtree();

CREATE TRIGGER habitats_validate_group
AFTER INSERT OR UPDATE OF habitat_group_id ON habitats
FOR EACH ROW EXECUTE FUNCTION validate_group_parent_not_in_subtree();
