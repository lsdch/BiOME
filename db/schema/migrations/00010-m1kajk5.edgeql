CREATE MIGRATION m1kajk5j65rxizvzue3wwq5fnehv2xica75s7u4a6i3v3cgpfipbxq
    ONTO m17kzjcbuahadpg2fgi3vkgmcszgpcck3lb4c3pt6uy77z57xelnoq
{
  CREATE FUNCTION events::sampling_code(s: events::Sampling, i: std::int32) ->  std::str USING (WITH
      date_suffix := 
          (SELECT
              std::assert_single((IF (s.event.performed_on.precision = date::DatePrecision.Unknown) THEN 'UNK' ELSE (IF (s.event.performed_on.precision = date::DatePrecision.Year) THEN <std::str>std::datetime_get(s.event.performed_on.date, 'year') ELSE (<std::str>std::datetime_get(s.event.performed_on.date, 'year') ++ std::str_pad_start(<std::str>std::datetime_get(s.event.performed_on.date, 'month'), 2, '0')))))
          )
      ,
      base := 
          std::assert_single(((s.event.site.code ++ '_') ++ date_suffix))
  SELECT
      std::assert_exists(std::assert_single((IF (i > 0) THEN ((base ++ '.') ++ <std::str>i) ELSE base)))
  );
  ALTER TYPE events::Sampling {
      ALTER PROPERTY code {
          DROP ANNOTATION default::example;
          DROP ANNOTATION std::description;
          DROP CONSTRAINT std::exclusive;
          DROP REWRITE
              INSERT ;
              DROP REWRITE
                  UPDATE ;
              };
          };
  ALTER TYPE samples::BioMaterial {
      ALTER PROPERTY code {
          DROP REWRITE
              INSERT ;
          };
      };
  ALTER TYPE samples::BioMaterial {
      ALTER PROPERTY code {
          DROP REWRITE
              UPDATE ;
          };
      };
  ALTER TYPE events::Sampling {
      DROP PROPERTY code;
  };
  ALTER TYPE events::Sampling {
      CREATE REQUIRED PROPERTY suffix_id: std::int32 {
          SET default := 1;
      };
  };
  ALTER TYPE events::Sampling {
      CREATE REQUIRED PROPERTY code := (events::sampling_code(__source__, .suffix_id));
  };
  ALTER TYPE samples::BioMaterial {
      ALTER PROPERTY code {
          CREATE REWRITE
              INSERT 
              USING (((.identification.taxon.code ++ '|') ++ .sampling.code));
      };
  };
  ALTER TYPE samples::BioMaterial {
      ALTER PROPERTY code {
          CREATE REWRITE
              UPDATE 
              USING (((.identification.taxon.code ++ '|') ++ .sampling.code));
      };
  };
  CREATE TYPE events::SamplingCodeIndex {
      CREATE REQUIRED LINK site: location::Site;
      CREATE REQUIRED PROPERTY count: std::int32;
      CREATE REQUIRED PROPERTY month: std::int32;
      CREATE REQUIRED PROPERTY year: std::int32;
      CREATE CONSTRAINT std::exclusive ON ((.site, .year, .month));
      CREATE INDEX ON ((.site, .year, .month));
  };
  ALTER TYPE events::Event {
      CREATE TRIGGER event_update
          AFTER UPDATE 
          FOR EACH 
              WHEN (((<std::int32>std::datetime_get(__new__.performed_on.date, 'year') != <std::int32>std::datetime_get(__old__.performed_on.date, 'year')) OR (<std::int32>std::datetime_get(__new__.performed_on.date, 'month') != <std::int32>std::datetime_get(__old__.performed_on.date, 'month'))))
          DO (INSERT
              events::SamplingCodeIndex
              {
                  site := __new__.site,
                  year := <std::int32>std::datetime_get(__new__.performed_on.date, 'year'),
                  month := <std::int32>std::datetime_get(__new__.performed_on.date, 'month'),
                  count := std::count(__new__.samplings)
              }UNLESS CONFLICT ON (.site, .year, .month) ELSE (UPDATE
              events::SamplingCodeIndex
          SET {
              count := (.count + <std::int32>std::count(__new__.samplings))
          }));
  };
  ALTER TYPE events::Sampling {
      DROP PROPERTY generated_code;
      ALTER PROPERTY suffix_id {
          CREATE REWRITE
              INSERT 
              USING (<std::int32>((SELECT
                  events::SamplingCodeIndex
              FILTER
                  (((.site = __subject__.event.site) AND (.year = <std::int32>std::datetime_get(__subject__.event.performed_on.date, 'year'))) AND (.month = <std::int32>std::datetime_get(__subject__.event.performed_on.date, 'month')))
              )).count);
          CREATE REWRITE
              UPDATE 
              USING (<std::int32>((SELECT
                  events::SamplingCodeIndex
              FILTER
                  (((.site = __subject__.event.site) AND (.year = <std::int32>std::datetime_get(__subject__.event.performed_on.date, 'year'))) AND (.month = <std::int32>std::datetime_get(__subject__.event.performed_on.date, 'month')))
              )).count);
      };
      CREATE TRIGGER index_insert
          AFTER INSERT 
          FOR EACH DO (INSERT
              events::SamplingCodeIndex
              {
                  site := __new__.event.site,
                  year := <std::int32>std::datetime_get(__new__.event.performed_on.date, 'year'),
                  month := <std::int32>std::datetime_get(__new__.event.performed_on.date, 'month'),
                  count := 1
              }UNLESS CONFLICT ON (.site, .year, .month) ELSE (UPDATE
              events::SamplingCodeIndex
          SET {
              count := (.count + 1)
          }));
  };
};
