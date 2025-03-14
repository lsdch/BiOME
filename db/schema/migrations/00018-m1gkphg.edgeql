CREATE MIGRATION m1gkphgalqcyoeudgwixbuybsqgqou73ud6kgvojvsy7y5tjpeokrq
    ONTO m1ospj5bhcqtnqm3pvzrhoummiq5jyk7dhfibxz46njagvhkqy3cba
{
  ALTER FUNCTION location::insert_site(data: std::json) USING (WITH
      coords := 
          location::coords_from_json((data)['coordinates'])
      ,
      country_code := 
          <std::str>std::json_get(data, 'country_code')
  INSERT
      location::Site
      {
          name := <std::str>(data)['name'],
          code := <std::str>(data)['code'],
          description := <std::str>std::json_get(data, 'description'),
          coordinates := coords,
          locality := <std::str>std::json_get(data, 'locality'),
          country := (SELECT
              (IF EXISTS (country_code) THEN location::find_country(country_code) ELSE location::site_coords_to_country(coords))
          ),
          altitude := <std::int32>std::json_get(data, 'altitude')
      });
};
