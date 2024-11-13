CREATE MIGRATION m1cyfkatidiapmm7kxfcpye4nvcpn6xf36exrafcecuv7osw7lsufa
    ONTO m1d7nxsburgqsiaesq3guxncqorelaxa76q7js6zhda3wr4iynd4na
{
  CREATE FUNCTION events::event_code(e: events::Event) ->  std::str USING (WITH
      date_suffix := 
          (SELECT
              std::assert_single((IF (e.performed_on.precision = date::DatePrecision.Unknown) THEN 'UNK' ELSE (IF (e.performed_on.precision = date::DatePrecision.Year) THEN <std::str>std::datetime_get(e.performed_on.date, 'year') ELSE (<std::str>std::datetime_get(e.performed_on.date, 'year') ++ std::str_pad_start(<std::str>std::datetime_get(e.performed_on.date, 'month'), 2, '0')))))
          )
  SELECT
      ((e.site.code ++ '_') ++ date_suffix)
  );
  ALTER TYPE samples::BioMaterial {
      ALTER PROPERTY code {
          DROP REWRITE
              INSERT ;
          };
      };
  ALTER TYPE samples::BioMaterial {
      ALTER PROPERTY code {
          CREATE REWRITE
              INSERT 
              USING (((.identification.taxon.code ++ '|') ++ events::event_code(.sampling.event)));
      };
  };
  ALTER TYPE samples::BioMaterial {
      ALTER PROPERTY code {
          DROP REWRITE
              UPDATE ;
          };
      };
  ALTER TYPE samples::BioMaterial {
      ALTER PROPERTY code {
          CREATE REWRITE
              UPDATE 
              USING (((.identification.taxon.code ++ '|') ++ events::event_code(.sampling.event)));
      };
  };
  ALTER TYPE events::Sampling {
      DROP PROPERTY code;
  };
};
