CREATE MIGRATION m16ybjkvw3prnvznuep5aw36f6gffcn24kbuf53jltu5t4je7lyvna
    ONTO m1kjfhbmuwa65j4byyzeato7bzon6ydqxjpl2kttrz45izfmwqddba
{
  ALTER TYPE default::Auditable {
      ALTER LINK meta {
          SET readonly := true;
          DROP CONSTRAINT std::exclusive;
      };
  };
};
