CREATE MIGRATION m1rbvgir64ynum5rwujqqgnq3mrpvtgjpteqhupd2seg7heyi7f4qa
    ONTO m1ep57s3srib5cagyluptd6e54twnhx34gwyexuu3zrijblk6joq6a
{
  ALTER TYPE location::Site {
      ALTER PROPERTY name {
          RESET OPTIONALITY;
      };
  };
  ALTER FUNCTION location::insert_site(data: std::json) USING (WITH
      coords := 
          location::coords_from_json((data)['coordinates'])
      ,
      country_code := 
          <std::str>std::json_get(data, 'country_code')
  INSERT
      location::Site
      {
          name := <std::str>std::json_get(data, 'name'),
          code := <std::str>(data)['code'],
          description := <std::str>std::json_get(data, 'description'),
          coordinates := coords,
          locality := <std::str>std::json_get(data, 'locality'),
          country := (SELECT
              (IF EXISTS (country_code) THEN location::find_country(country_code) ELSE location::site_coords_to_country(coords))
          ),
          altitude := <std::int32>std::json_get(data, 'altitude')
      });
  CREATE FUNCTION taxonomy::is_child(taxon: taxonomy::Taxon, ancestor: taxonomy::Taxon) ->  std::bool USING (std::assert_exists((IF (ancestor.rank = taxonomy::Rank.Kingdom) THEN (taxon.kingdom = ancestor) ELSE (IF (ancestor.rank = taxonomy::Rank.Phylum) THEN (taxon.phylum = ancestor) ELSE (IF (ancestor.rank = taxonomy::Rank.Class) THEN (taxon.class = ancestor) ELSE (IF (ancestor.rank = taxonomy::Rank.Order) THEN (taxon.order = ancestor) ELSE (IF (ancestor.rank = taxonomy::Rank.Family) THEN (taxon.family = ancestor) ELSE (IF (ancestor.rank = taxonomy::Rank.Genus) THEN (taxon.genus = ancestor) ELSE (IF (ancestor.rank = taxonomy::Rank.Species) THEN (taxon.species = ancestor) ELSE false)))))))));
};
