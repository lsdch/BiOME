CREATE MIGRATION m1jri6xhsc227mifrbjrq6kvtxa65tglz6pticb2fa6maocqdoivba
    ONTO m135y5kqptyislr7zwamflhfcws7gwwdh4voesedgsvdvggjot4jna
{
  ALTER TYPE taxonomy::SynonymGroup {
      CREATE REQUIRED LINK accepted: taxonomy::Taxon {
          SET REQUIRED USING (<taxonomy::Taxon>{});
          CREATE CONSTRAINT std::exclusive;
      };
      CREATE REQUIRED MULTI LINK synonyms := ((.taxa UNION .accepted));
  };
  ALTER TYPE taxonomy::Taxon {
      DROP LINK synonym_group;
  };
  ALTER TYPE taxonomy::Taxon {
      CREATE LINK synonym_group := ((.<accepted[IS taxonomy::SynonymGroup] UNION .<taxa[IS taxonomy::SynonymGroup]));
      ALTER LINK synonyms {
          USING ((.synonym_group.synonyms EXCEPT __source__));
      };
  };
};
