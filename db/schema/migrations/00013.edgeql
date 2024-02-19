CREATE MIGRATION m1ur3upk4dbze2w4jahujyhe6pfhsdqdsz2yyf3gnwhajh353omnja
    ONTO m15tsemmruo62xzzgb7jczincphvpnfbsvznardetpomlezd2ls2da
{
  ALTER TYPE taxonomy::Taxon {
      CREATE PROPERTY comment: std::str;
  };
};
