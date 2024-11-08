CREATE MIGRATION m13nmsmvfnx6bnalbpk245ea2fupkbamdua7csqyfwzbz7ij73xa4a
    ONTO m1t2ij46mbc7qowsixuauhblvwolammlhjvidhgjkeknj56epxmqyq
{
  ALTER TYPE events::Spotting {
      ALTER LINK event {
          ON TARGET DELETE DELETE SOURCE;
          CREATE CONSTRAINT std::exclusive;
          SET REQUIRED;
      };
  };
  ALTER TYPE events::Event {
      CREATE LINK spotting := (.<event[IS events::Spotting]);
  };
  ALTER TYPE events::Event {
      DROP LINK spottings;
  };
  ALTER TYPE events::Spotting {
      ALTER LINK event {
          SET OWNED;
          SET TYPE events::Event;
      };
  };
};
