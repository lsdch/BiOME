CREATE MIGRATION m1pneippm6lkmhrkjdzsrkbcjw7agimxre4nkxp6aixzpnj5ycpswa
    ONTO m1y5ynh7hph447hhf2dmz6lpi5lkaj2pnueuxk3qqr3pvglilephxq
{
  ALTER TYPE event::EventDataset {
      DROP LINK assembled_by;
  };
  ALTER TYPE event::EventDataset {
      CREATE REQUIRED MULTI LINK maintainers: people::Person {
          ON TARGET DELETE RESTRICT;
          SET REQUIRED USING (<people::Person>{});
      };
  };
};
