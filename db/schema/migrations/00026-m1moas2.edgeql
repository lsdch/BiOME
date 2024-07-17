CREATE MIGRATION m1moas2expbrkao543rjlehyh5i7sc7rw6zcgdkwxu2eptwkghxgca
    ONTO m1d3ovccfcoctdodj7fmmjpqixjqezdgazbhythtsjveu7xyrd26oa
{
  ALTER TYPE taxonomy::Taxon {
      CREATE CONSTRAINT std::expression ON ((std::len(std::str_split(.name, ' ')) >= 3)) EXCEPT ((.rank != taxonomy::Rank.Subspecies)) {
          SET errmessage := 'A subspecies name must include at least 2 whitespaces.';
      };
  };
  ALTER TYPE taxonomy::Taxon {
      CREATE CONSTRAINT std::expression ON ((std::len(std::str_split(.name, ' ')) >= 2)) EXCEPT ((.rank != taxonomy::Rank.Species)) {
          SET errmessage := 'A species name must include a whitespace.';
      };
  };
  ALTER TYPE taxonomy::Taxon {
      DROP CONSTRAINT std::expression ON ((std::len(std::str_split(.name, ' ')) = 3)) EXCEPT ((.rank != taxonomy::Rank.Subspecies));
      ALTER CONSTRAINT std::expression ON (NOT (std::contains(.name, ' '))) EXCEPT (((.rank = taxonomy::Rank.Species) OR (.rank = taxonomy::Rank.Subspecies))) {
          SET errmessage := 'Taxon names with rank higher than species may not include a whitespace.';
      };
  };
  ALTER TYPE taxonomy::Taxon {
      DROP CONSTRAINT std::expression ON ((std::len(std::str_split(.name, ' ')) = 2)) EXCEPT ((.rank != taxonomy::Rank.Species));
  };
};
