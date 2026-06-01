-- View to easily query synonyms and their accepted taxa
CREATE VIEW taxon_synonyms AS
SELECT accepted.id AS accepted_taxon_id,
	accepted.name AS accepted_taxon_name,
	synonym.id AS synonym_taxon_id,
	synonym.name AS synonym_taxon_name
FROM taxa accepted
	JOIN taxa synonym ON synonym.accepted_taxon_id = accepted.id;

-- Taxonomic hierarchy view for efficient querying of ancestor-descendant relationships
CREATE VIEW taxon_hierarchy AS WITH RECURSIVE tree AS (
	SELECT t.id,
		t.parent_id,
		t.name,
		t.authorship,
		t.rank,
		t.status,
		t.comments,
		t.id AS root_id,
		ARRAY [t.id] AS path,
		0 AS depth
	FROM taxa t
	WHERE t.parent_id IS NULL
	UNION ALL
	SELECT c.id,
		c.parent_id,
		c.name,
		c.authorship,
		c.rank,
		c.status,
		c.comments,
		tree.root_id,
		tree.path || c.id,
		tree.depth + 1
	FROM taxa c
		JOIN tree ON c.parent_id = tree.id
)
SELECT *
FROM tree;