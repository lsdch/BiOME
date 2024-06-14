CREATE MIGRATION m1vev6tejtldvpak7nk7itzueju3bss34as6gy3zlckmvavjympfoq
    ONTO m1xlljiu42rsst325lg6znmcq5ofu6ktsabx72r3crqq5rsgkvdpia
{
  ALTER TYPE location::Locality {
      DROP LINK sites;
      DROP PROPERTY code;
  };
  ALTER TYPE location::Site {
      CREATE MULTI LINK datasets := (.<sites[IS location::SiteDataset]);
      DROP LINK habitat_tags;
  };
  ALTER TYPE location::Site {
      ALTER PROPERTY altitudeRange {
          RENAME TO altitude;
      };
  };
};
