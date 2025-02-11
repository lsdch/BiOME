CREATE MIGRATION m1uicsjoj2sfqrubiqwsyhoxyfsl6wh2ggr7flrhfyrlncpo6seyga
    ONTO m1z2chv3zjiy7grfmpca3x5bb3vqtbdicahvkc27yhzmsw4ob5urqq
{
  ALTER FUNCTION location::insert_site(data: std::json) USING (INSERT
      location::Site
      {
          name := <std::str>(data)['name'],
          code := <std::str>(data)['code'],
          description := <std::str>std::json_get(data, 'description'),
          coordinates := location::coords_from_json((data)['coordinates']),
          locality := <std::str>std::json_get(data, 'locality'),
          country := location::find_country(<std::str>std::json_get(data, 'country_code')),
          altitude := <std::int32>std::json_get(data, 'altitude')
      });
};
