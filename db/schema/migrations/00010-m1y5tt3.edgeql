CREATE MIGRATION m1y5tt37vrqrgw37uztchtyrbyxu7re7e426y2hxt5vhdsqxsmbr6q
    ONTO m1azts2axrtmmlo2b3fkfqaqewyemgqfrznstpkijyseqg3chttapq
{
  ALTER TYPE occurrence::Occurrence {
      ALTER PROPERTY code {
          SET default := (<std::str>std::sequence_next(INTROSPECT occurrence::DefaultOccurrenceCode));
          CREATE REWRITE
              INSERT 
              USING ((IF GLOBAL occurrence::block_code_update THEN .code ELSE (((taxonomy::taxon_code(.identification.taxon) ++ '[') ++ .sampling.code) ++ ']')));
          CREATE REWRITE
              UPDATE 
              USING ((IF GLOBAL occurrence::block_code_update THEN .code ELSE (((taxonomy::taxon_code(.identification.taxon) ++ '[') ++ .sampling.code) ++ ']')));
      };
  };
};
