CREATE MIGRATION m1z3yrk4amxwng2hqa2nlkg2ibx6juobzz577appsckeo2bjlp2xdq
    ONTO m1kfabjj3mnl3zjkc5m2x2waiw2gzcz4shwh3ydb7msc5howge3o7q
{
  ALTER FUNCTION taxonomy::taxon_code(taxon: taxonomy::Taxon) USING (std::re_replace('[^()a-zA-Z]', '_', std::str_replace(taxon.name, 'sp. ', 'sp.'), flags := 'ig'));
};
