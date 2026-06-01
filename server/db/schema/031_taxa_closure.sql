CREATE TABLE taxa_closure (
    ancestor_id UUID NOT NULL REFERENCES taxa(id) ON DELETE CASCADE,
    descendant_id UUID NOT NULL REFERENCES taxa(id) ON DELETE CASCADE,
    depth INTEGER NOT NULL,
    -- 0 = self, 1 = parent, etc.
    PRIMARY KEY (ancestor_id, descendant_id)
);

CREATE INDEX taxa_closure_ancestor_idx ON taxa_closure (ancestor_id);

CREATE INDEX taxa_closure_descendant_idx ON taxa_closure (descendant_id);

CREATE INDEX taxa_closure_depth_idx ON taxa_closure (depth);

CREATE INDEX taxa_closure_cycle_check ON taxa_closure (ancestor_id, descendant_id);

---------------------------------------------------------------------
-- Insert trigger to maintain the closure table
---------------------------------------------------------------------
CREATE OR REPLACE FUNCTION taxa_closure_insert_fn() RETURNS TRIGGER AS $$ BEGIN --
-- 1. self relation
INSERT INTO taxa_closure (ancestor_id, descendant_id, depth)
VALUES (NEW.id, NEW.id, 0);

-- 2. inherit ancestors from parent
IF NEW.parent_id IS NOT NULL THEN
INSERT INTO taxa_closure (ancestor_id, descendant_id, depth)
SELECT ancestor_id,
    NEW.id,
    depth + 1
FROM taxa_closure
WHERE descendant_id = NEW.parent_id;
END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_taxa_closure_insert
AFTER
INSERT ON taxa FOR EACH ROW EXECUTE FUNCTION taxa_closure_insert_fn();

---------------------------------------------------------------------
-- Delete trigger to maintain the closure table
---------------------------------------------------------------------
CREATE OR REPLACE FUNCTION taxa_closure_delete_fn() RETURNS TRIGGER AS $$ BEGIN
DELETE FROM taxa_closure
WHERE descendant_id IN (
        SELECT descendant_id
        FROM taxa_closure
        WHERE ancestor_id = OLD.id
    );

    RETURN OLD;
END;
$$ LANGUAGE plpgsql;
CREATE TRIGGER trg_taxa_closure_delete
AFTER DELETE ON taxa FOR EACH ROW EXECUTE FUNCTION taxa_closure_delete_fn();

---------------------------------------------------------------------
-- Update trigger to maintain the closure table
---------------------------------------------------------------------
CREATE OR REPLACE FUNCTION taxa_closure_update_fn() RETURNS TRIGGER AS $$ BEGIN -- si parent inchangé, rien à faire
    IF NEW.parent_id IS NOT DISTINCT
FROM OLD.parent_id THEN RETURN NEW;
END IF;

-- delete old paths
DELETE FROM taxa_closure
WHERE descendant_id IN (
        SELECT descendant_id
        FROM taxa_closure
        WHERE ancestor_id = OLD.id
    );

-- rebuild self link
INSERT INTO taxa_closure (ancestor_id, descendant_id, depth)
VALUES (NEW.id, NEW.id, 0);

-- rebuild paths from new parent
IF NEW.parent_id IS NOT NULL THEN
INSERT INTO taxa_closure (ancestor_id, descendant_id, depth)
SELECT ancestor_id,
    child.descendant_id,
    parent.depth + child.depth + 1
FROM taxa_closure parent
    JOIN taxa_closure child ON child.ancestor_id = OLD.id
WHERE parent.descendant_id = NEW.parent_id;
END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_taxa_closure_update
AFTER
UPDATE OF parent_id ON taxa FOR EACH ROW EXECUTE FUNCTION taxa_closure_update_fn();