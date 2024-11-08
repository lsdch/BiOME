CREATE MIGRATION m1wcpeg6y4xz6lbdolthika34bfm4xldvk4gkxl7u7jtgz7cr4naia
    ONTO m1lzidennerhb2nvvhipwh4ecvwhdgimsdx45v7vcpiqizl6nyshja
{
  ALTER TYPE events::Event {
      CREATE LINK dataset: location::SiteDataset;
  };
  ALTER TYPE events::Sampling {
      ALTER PROPERTY code {
          SET default := '';
          ALTER ANNOTATION default::example := 'SOMESITE_202301 ; SOMESITE_202301.1';
          ALTER ANNOTATION std::description := 'Format : SITE_YEARMONTH.NUMBER. The NUMBER suffix is not appended if the site and month tuple is unique.';
          DROP ANNOTATION std::title;
      };
      ALTER PROPERTY generated_code {
          USING (SELECT
              (((.event.site.code ++ '_') ++ <std::str>std::datetime_get(.event.performed_on.date, 'year')) ++ std::str_pad_start(<std::str>std::datetime_get(.event.performed_on.date, 'month'), 2, '0'))
          );
          SET REQUIRED;
      };
      ALTER PROPERTY sampling_duration {
          SET TYPE std::int32 USING (<std::int32>std::duration_get(.sampling_duration, 'minute'));
      };
  };
  ALTER TYPE location::Site {
      ALTER LINK datasets {
          RESET EXPRESSION;
          RESET EXPRESSION;
          ON SOURCE DELETE ALLOW;
          ON TARGET DELETE ALLOW;
          RESET OPTIONALITY;
          SET TYPE location::SiteDataset;
      };
      ALTER PROPERTY coordinates {
          CREATE REWRITE
              INSERT 
              USING ((
                  precision := __subject__.coordinates.precision,
                  latitude := <std::float32>std::round(<std::decimal>__subject__.coordinates.latitude, 5),
                  longitude := <std::float32>std::round(<std::decimal>__subject__.coordinates.longitude, 5)
              ));
          CREATE REWRITE
              UPDATE 
              USING ((
                  precision := __subject__.coordinates.precision,
                  latitude := <std::float32>std::round(<std::decimal>__subject__.coordinates.latitude, 5),
                  longitude := <std::float32>std::round(<std::decimal>__subject__.coordinates.longitude, 5)
              ));
      };
  };
  ALTER TYPE location::SiteDataset {
      ALTER LINK sites {
          USING (.<datasets[IS location::Site]);
          RESET ON TARGET DELETE;
      };
  };
};
