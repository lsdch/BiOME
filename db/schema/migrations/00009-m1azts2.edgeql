CREATE MIGRATION m1azts2axrtmmlo2b3fkfqaqewyemgqfrznstpkijyseqg3chttapq
    ONTO m1y7dc5hxzyu74a62utzt3hl6evgnlfxmv3qyrxwbll4recrvjcsqq
{
  CREATE SCALAR TYPE occurrence::DefaultOccurrenceCode EXTENDING std::sequence;
  CREATE GLOBAL occurrence::block_code_update -> std::bool {
      SET default := false;
      CREATE ANNOTATION std::description := 'When true, prevents updates to the occurrence code. This is used to enforce immutability of occurrence codes after creation.';
  };
};
