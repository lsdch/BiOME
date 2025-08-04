CREATE MIGRATION m1aik4lkvpgejpmxfdyrerxz3lchi7ovfcajv44vyycl3yqj4ypnna
    ONTO m1vhuzbfmksam34as2anx5jomwdjlmzo5wf5uav4tzaqjnssefkgzq
{
  ALTER TYPE datasets::Dataset {
      CREATE LINK publication: references::Article {
          ON SOURCE DELETE ALLOW;
      };
  };
};
