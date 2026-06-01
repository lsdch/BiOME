-- Habitat groups and habitat tags
CREATE TABLE habitat_groups (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
	label CITEXT NOT NULL,
	CONSTRAINT habitat_group_label_unique UNIQUE (label),
	description TEXT,
	exclusive_elements BOOLEAN NOT NULL DEFAULT true,
	CONSTRAINT habitat_group_label_not_empty CHECK (btrim(label) <> '')
);

CREATE UNIQUE INDEX idx_habitat_group_label_uq ON habitat_groups (label);

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
DECLARE is_invalid BOOLEAN;
BEGIN IF NEW.parent_habitat_id IS NULL THEN RETURN NEW;
END IF;

    WITH RECURSIVE descendants AS (
	-- tous les habitats du groupe
	SELECT h.id
	FROM habitats h
	WHERE h.habitat_group_id = NEW.id
	UNION ALL
	-- descendance récursive
	SELECT child.id
	FROM habitats child
		JOIN descendants d ON child.parent_id = d.id
)
SELECT TRUE INTO is_invalid
FROM descendants
WHERE id = NEW.parent_habitat_id
LIMIT 1;

    IF is_invalid THEN RAISE EXCEPTION 'Invalid parent_habitat_id: cannot reference a descendant habitat of the same group' USING ERRCODE = 'HB001';
END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;