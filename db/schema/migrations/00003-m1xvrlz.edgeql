CREATE MIGRATION m1xvrlz4my3tzflk3k6vf2bnb2e2k5zzzljxslvmgpyjsder5s6aaq
    ONTO m1b2dodwjgwhy7vfu5lsvvpc7njb35xtvtuuwk7jgmdezjegr4niza
{
  ALTER TYPE location::SiteDataset {
      CREATE REQUIRED PROPERTY slug: std::str {
          SET REQUIRED USING (<std::str>{});
          CREATE CONSTRAINT std::exclusive;
      };
  };
};
