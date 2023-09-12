CREATE MIGRATION m1nwy6dl5u7s4k6nb6pmn5y4lpijy36emjll6l3xjmy2ufuxwvxaza
    ONTO m1pneippm6lkmhrkjdzsrkbcjw7agimxre4nkxp6aixzpnj5ycpswa
{
  CREATE TYPE location::SiteDataset EXTENDING default::Auditable {
      CREATE REQUIRED MULTI LINK maintainers: people::Person;
      CREATE MULTI LINK sites: location::Site;
      CREATE PROPERTY description: std::str;
      CREATE REQUIRED PROPERTY label: std::str {
          CREATE CONSTRAINT std::max_len_value(40);
          CREATE CONSTRAINT std::min_len_value(4);
      };
  };
};
