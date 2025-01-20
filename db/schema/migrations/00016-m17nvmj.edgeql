CREATE MIGRATION m17nvmjy53v4uqscq3jp3xnwvllww6jnpn3wit6gpcplwq42xwr4uq
    ONTO m1x74vqa5iyqskd3eahh6ffanjidjadcaetzzg67hnhfnckx55jgqq
{
  CREATE TYPE admin::GeoapifyUsage {
      CREATE REQUIRED PROPERTY date: std::cal::local_date {
          CREATE CONSTRAINT std::exclusive;
      };
      CREATE REQUIRED PROPERTY requests: std::int32;
  };
};
