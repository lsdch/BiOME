CREATE MIGRATION m14xzezsvdeffja6ewwjv4ngugttyli6oifdiihuqb7dzv3cwueoba
    ONTO m1jri6xhsc227mifrbjrq6kvtxa65tglz6pticb2fa6maocqdoivba
{
  ALTER TYPE taxonomy::SynonymGroup {
      ALTER LINK accepted {
          DROP CONSTRAINT std::exclusive;
      };
  };
  ALTER TYPE taxonomy::Taxon {
      DROP LINK synonyms;
  };
  ALTER TYPE taxonomy::Taxon {
      DROP LINK synonym_group;
  };
  ALTER TYPE taxonomy::Taxon {
      CREATE LINK synonym_group: taxonomy::SynonymGroup {
          ON SOURCE DELETE DELETE TARGET IF ORPHAN;
          ON TARGET DELETE ALLOW;
      };
  };
  ALTER TYPE taxonomy::SynonymGroup {
      ALTER LINK synonyms {
          USING (std::assert_exists(.<synonym_group[IS taxonomy::Taxon]));
      };
      DROP LINK accepted;
      DROP LINK taxa;
  };
  ALTER TYPE taxonomy::Taxon {
      CREATE MULTI LINK synonyms := ((.synonym_group.synonyms EXCEPT __source__));
  };
};
