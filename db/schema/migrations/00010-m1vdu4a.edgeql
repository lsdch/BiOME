CREATE MIGRATION m1vdu4awpotx3hw3imqe3dubgebh5fnvcibyd36m4m4qpwewdtmp7a
    ONTO m1xeampwczkkzmxspppuhrczxk6jybxqcg72zyn7tekvoppmxi6tva
{
  ALTER TYPE location::Site {
      ALTER PROPERTY name {
          DROP CONSTRAINT std::exclusive;
      };
  };
};
