CREATE MIGRATION m1c52quj6kiab4vp6xl5efdwlmrssvcxpzwq6mphvzg2zitcnj2gta
    ONTO m1kajk5j65rxizvzue3wwq5fnehv2xica75s7u4a6i3v3cgpfipbxq
{
  ALTER TYPE events::Sampling {
      ALTER PROPERTY suffix_id {
          DROP REWRITE
              INSERT ;
          };
      };
  ALTER TYPE events::Sampling {
      ALTER PROPERTY suffix_id {
          CREATE REWRITE
              INSERT 
              USING (<std::int32>(SELECT
                  (((SELECT
                      events::SamplingCodeIndex
                  FILTER
                      (((.site = __subject__.event.site) AND (.year = <std::int32>std::datetime_get(__subject__.event.performed_on.date, 'year'))) AND (.month = <std::int32>std::datetime_get(__subject__.event.performed_on.date, 'month')))
                  )).count ?? 1)
              ));
      };
  };
  ALTER TYPE events::Sampling {
      ALTER PROPERTY suffix_id {
          DROP REWRITE
              UPDATE ;
          };
      };
  ALTER TYPE events::Sampling {
      ALTER PROPERTY suffix_id {
          CREATE REWRITE
              UPDATE 
              USING (<std::int32>(SELECT
                  (((SELECT
                      events::SamplingCodeIndex
                  FILTER
                      (((.site = __subject__.event.site) AND (.year = <std::int32>std::datetime_get(__subject__.event.performed_on.date, 'year'))) AND (.month = <std::int32>std::datetime_get(__subject__.event.performed_on.date, 'month')))
                  )).count ?? 1)
              ));
      };
  };
};
