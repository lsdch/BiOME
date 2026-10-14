CREATE TABLE habitats (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    label CITEXT NOT NULL,
    facet CITEXT,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT habitat_label_not_empty CHECK (btrim(label) <> ''),
    CONSTRAINT habitat_facet_not_empty CHECK (
        facet IS NULL
        OR btrim(facet) <> ''
    ),
    CONSTRAINT habitat_description_length CHECK (char_length(coalesce(description, '')) <= 4000),
    CONSTRAINT uq_habitat_label UNIQUE (label)
);

CREATE TABLE samplings_habitats (
    sampling_id UUID NOT NULL REFERENCES samplings (id) ON DELETE CASCADE,
    habitat_id UUID NOT NULL REFERENCES habitats (id) ON DELETE CASCADE,
    PRIMARY KEY (sampling_id, habitat_id)
);
CREATE INDEX samplings_habitats_habitat_idx ON samplings_habitats (habitat_id);