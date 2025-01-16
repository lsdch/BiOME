CREATE MIGRATION m1mhsnbai2qljuo5j2ks4amkpli5wsjm25muhdikvivjb2jt6v5pza
    ONTO m1vdu4awpotx3hw3imqe3dubgebh5fnvcibyd36m4m4qpwewdtmp7a
{
  ALTER TYPE location::Site {
      ALTER PROPERTY code {
          DROP CONSTRAINT std::min_len_value(4);
      };
  };
  ALTER TYPE location::Site {
      ALTER PROPERTY code {
          CREATE CONSTRAINT std::min_len_value(3);
      };
  };
};
