CREATE MIGRATION m13xy2wljavngnlipsmytisyzyr2acgfwqd4czg2subfvylqcqvsxq
    ONTO m15eeizjt2q7igd2sjkwy5oz6jicc5bmpggm7p4kfapn62bdblm4uq
{
  ALTER TYPE event::PlannedSampling {
      ALTER LINK target_taxa {
          SET REQUIRED USING (<taxonomy::Taxon>{});
      };
      DROP PROPERTY comments;
  };
  ALTER TYPE taxonomy::Taxon {
      ALTER LINK parent {
          ON TARGET DELETE DELETE SOURCE;
      };
  };
};
