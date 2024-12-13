CREATE MIGRATION m1uqp4r4s5hj6oa5fsrjymkhjmkhh2fw6sndccn3cwyoevxsgqgexa
    ONTO m1gkyaucljimawntjiae6426sjzubgnes5m5jfvoiwiqu6lld27grq
{
  ALTER TYPE occurrence::BioMaterial {
      CREATE PROPERTY code_history: array<tuple<code: std::str, time: std::datetime>> {
          SET readonly := true;
          CREATE REWRITE
              UPDATE 
              USING ((IF (__old__.code != .code) THEN (__old__.code_history ++ [(
                  code := __old__.code,
                  time := std::datetime_of_statement()
              )]) ELSE .code_history));
      };
  };
  DROP FUNCTION events::event_code(e: events::Event);
  ALTER TYPE events::Event {
      CREATE REQUIRED PROPERTY code := (WITH
          date := 
              .performed_on.date
          ,
          precision := 
              .performed_on.precision
      SELECT
          ((.site.code ++ '|') ++ (IF (precision = date::DatePrecision.Unknown) THEN 'undated' ELSE (IF (precision = date::DatePrecision.Year) THEN <std::str>std::datetime_get(date, 'year') ELSE ((<std::str>std::datetime_get(date, 'year') ++ '-') ++ <std::str>std::datetime_get(date, 'month')))))
      );
  };
  CREATE SCALAR TYPE events::SamplingNumber EXTENDING std::sequence;
  ALTER TYPE events::Sampling {
      CREATE REQUIRED PROPERTY number: events::SamplingNumber {
          SET readonly := true;
          SET REQUIRED USING (<events::SamplingNumber>{});
      };
      CREATE REQUIRED SINGLE PROPERTY code := (WITH
          id := 
              .id
          ,
          event := 
              .event
          ,
          sisters := 
              (SELECT
                  DETACHED events::Sampling
              FILTER
                  (.event.code = event.code)
              ORDER BY
                  .number ASC
              )
          ,
          rank := 
              (SELECT
                  (std::assert_single(std::enumerate(sisters) FILTER
                      (.1.id = id)
                  )).0
              )
          ,
          suffix := 
              (IF (std::count(sisters) > 1) THEN ('.' ++ <std::str>(rank + 1)) ELSE '')
      SELECT
          std::assert_single(std::assert_exists((event.code ++ suffix)))
      );
  };
  ALTER TYPE occurrence::BioMaterial {
      ALTER PROPERTY code {
          ALTER ANNOTATION std::description := "Format like 'taxon_short_code[sampling_code]'";
          CREATE REWRITE
              UPDATE 
              USING ((IF __specified__.code THEN .code ELSE (((.identification.taxon.code ++ '[') ++ .sampling.code) ++ ']')));
      };
  };
  ALTER TYPE seq::ExternalSequence {
      ALTER PROPERTY code {
          ALTER CONSTRAINT std::exclusive {
              SET OWNED;
          };
          SET OWNED;
          SET REQUIRED;
          SET TYPE std::str;
      };
  };
  ALTER TYPE seq::SeqReference {
      CREATE REQUIRED PROPERTY code := (((.db.code ++ ':') ++ .accession));
      CREATE REQUIRED PROPERTY is_origin: std::bool {
          SET default := false;
      };
  };
  ALTER TYPE seq::ExternalSequence {
      ALTER PROPERTY code {
          CREATE REWRITE
              UPDATE 
              USING (WITH
                  suffix := 
                      (IF (.origin = seq::ExtSeqOrigin.Lab) THEN 'lab' ELSE (IF (.origin = seq::ExtSeqOrigin.PersCom) THEN 'perscom' ELSE (WITH
                          sources := 
                              ((SELECT
                                  .referenced_in
                              FILTER
                                  .is_origin
                              )).code
                      SELECT
                          std::array_join(std::array_agg(sources), '|')
                      )))
              SELECT
                  ((((((.identification.taxon.code ++ '[') ++ .sampling.code) ++ ']') ++ .specimen_identifier) ++ '|') ++ suffix)
              );
      };
  };
};
