CREATE MIGRATION m1ep57s3srib5cagyluptd6e54twnhx34gwyexuu3zrijblk6joq6a
    ONTO m15ohqkyku6ic33wyw4gyxcfg3crblpcedoprujokq67bwvf4ump3q
{
  ALTER TYPE datasets::ResearchProgram {
      CREATE MULTI LINK datasets: datasets::Dataset {
          ON SOURCE DELETE ALLOW;
          ON TARGET DELETE ALLOW;
      };
  };
  ALTER TYPE events::Event {
      DROP LINK programs;
  };
  DROP TYPE events::Program;
};
