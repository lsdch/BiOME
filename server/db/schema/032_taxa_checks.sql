-- Prevent taxon <-> parent cycles
CREATE OR REPLACE FUNCTION check_taxa_no_cycle () RETURNS trigger LANGUAGE plpgsql AS $$ BEGIN -- no parent => ok
    IF NEW.parent_id IS NULL THEN RETURN NEW;
END IF;

-- self-parenting
IF NEW.parent_id = NEW.id THEN RAISE EXCEPTION 'Taxon cannot be parent of itself' USING ERRCODE = 'TAX01',
DETAIL = format(
    'Taxon %s cannot be parent of itself',
    NEW.scientific_name
);
END IF;

-- synonym self-check 
IF NEW.accepted_taxon_id IS NOT NULL
AND NEW.accepted_taxon_id = NEW.id THEN RAISE EXCEPTION 'Taxon cannot be synonym of itself' USING ERRCODE = 'TAX02',
DETAIL = format(
    'Taxon %s cannot be synonym of itself',
    NEW.scientific_name
);
END IF;

-- CYCLE CHECK via closure table
IF EXISTS (
    SELECT 1
    FROM taxa_closure
    WHERE ancestor_id = NEW.id
        AND descendant_id = NEW.parent_id
) THEN RAISE EXCEPTION 'Cycle detected in taxon hierarchy' USING ERRCODE = 'TAX01';
END IF;

    RETURN NEW;
END;
$$;
CREATE TRIGGER taxa_no_cycle_trigger BEFORE
INSERT
    OR
UPDATE ON taxa FOR EACH ROW EXECUTE FUNCTION check_taxa_no_cycle ();