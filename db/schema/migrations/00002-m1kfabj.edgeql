CREATE MIGRATION m1kfabjj3mnl3zjkc5m2x2waiw2gzcz4shwh3ydb7msc5howge3o7q
    ONTO m1blq24e2cv7iyrttknc5t3fswzyekif5zycyyorwkkqbichvlulxa
{
  ALTER FUNCTION taxonomy::taxon_code(taxon: taxonomy::Taxon) USING (std::re_replace('[^()a-zA-Z]', '_', taxon.name));
};
