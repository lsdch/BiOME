CREATE MIGRATION m1fovhcwjrs6rznqorfxdso5weft3nr3ntnbqrrjviahz2vpybnziq
    ONTO m1dood75z46qkztfemwnx6ibzvowqxvlj3y53rrnafgjyct2lrzmsa
{
  DROP FUNCTION location::insert_site(data: std::json, infer_country: OPTIONAL std::bool);
  CREATE FUNCTION location::insert_site(data: std::json) ->  location::Site USING (WITH
      coords := 
          location::coords_from_json((data)['coordinates'])
  INSERT
      location::Site
      {
          name := <std::str>(data)['name'],
          code := <std::str>(data)['code'],
          description := <std::str>std::json_get(data, 'description'),
          coordinates := coords,
          locality := <std::str>std::json_get(data, 'locality'),
          country := (IF EXISTS (std::json_get(data, 'country_code')) THEN location::find_country(<std::str>std::json_get(data, 'country_code')) ELSE location::site_coords_to_country(coords)),
          altitude := <std::int32>std::json_get(data, 'altitude')
      });
};
