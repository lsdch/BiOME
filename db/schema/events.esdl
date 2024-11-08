module events {

  type Program extending default::Auditable {
    required label: str { constraint exclusive };
    required code: str { constraint exclusive };

    required multi managers: people::Person;
    multi funding_agencies: people::Institution;

    start_year: int32 {
      constraint min_value(1900);
    };
    end_year: int32;
    constraint expression on (.start_year <= .end_year);

    description: str;
  }

  type Event extending default::Auditable {
    required site: location::Site {
      on source delete allow;
      on target delete delete source;
    };

    required multi performed_by: people::Person;
    required performed_on: tuple<date: datetime, precision: date::DatePrecision> {
      constraint date::required_unless_unknown(__subject__.date, __subject__.precision);
      rewrite insert, update using (
        (
          date := date::rewrite_date(__subject__.performed_on),
          precision := __subject__.performed_on.precision
        )
      )
    };
    multi programs: Program;

    dataset: location::SiteDataset;

    spotting := .<event[is Spotting];
    multi abiotic_measurements := .<event[is AbioticMeasurement];
    multi samplings := .<event[is Sampling];

    trigger event_update after update for each
    when(
      <int32>datetime_get(__new__.performed_on.date, 'year') !=
      <int32>datetime_get(__old__.performed_on.date, 'year') or
      <int32>datetime_get(__new__.performed_on.date, 'month') !=
      <int32>datetime_get(__old__.performed_on.date, 'month')
    ) do (
        insert SamplingCodeIndex {
        site := __new__.site,
        year := <int32>datetime_get(__new__.performed_on.date, 'year'),
        month := <int32>datetime_get(__new__.performed_on.date, 'month'),
        count := count(__new__.samplings)
      } unless conflict on ((.site, .year, .month)) else (
        update SamplingCodeIndex set {
          count := .count + <int32>count(__new__.samplings)
        }
      )
    );
  }

  # Several actions may have been performed during an event
  #
  # - Spotting: visiting a site without taking any samples, e.g. in preparation of future sampling
  # - Abiotic measurement: measuring environmental variables at the site
  # - Sampling: sampling for the presence of one or several target taxa
  abstract type Action extending default::Auditable {
    required event: Event {
      on target delete delete source;
    };
  }

  type Spotting extending Action {
    overloaded required event: Event {
      on target delete delete source;
      constraint exclusive;
    }
    multi target_taxa: taxonomy::Taxon;
    comments: str;
  }

  type AbioticParameter extending default::Vocabulary, default::Auditable {
    required unit: str;
  };

  type AbioticMeasurement extending Action {
    required param: AbioticParameter;
    required value: float32;
    constraint exclusive on ((.event, .param))
  }

  scalar type SamplingTarget extending enum<Community, Unknown, Taxa>;

  type SamplingMethod extending default::Vocabulary, default::Auditable;

  # function sampling_code(e: Event, i: int32) -> str using (
  #   with
  #     number_suffix := (
  #       <int32>((select SamplingCodeIndex
  #         filter .site = e.site
  #         and .year = <int32>datetime_get(e.performed_on.date, 'year')
  #         and .month = <int32>datetime_get(e.performed_on.date, 'month')
  #       ).count ?? 1) + i
  #     ),
  #     date_suffix := (select assert_single(
  #       if(e.performed_on.precision = date::DatePrecision.Unknown)
  #       then "UNK"
  #       else if (e.performed_on.precision = date::DatePrecision.Year)
  #       then <str>datetime_get(e.performed_on.date, 'year')
  #       else (
  #         <str>datetime_get(e.performed_on.date, 'year') ++
  #         str_pad_start(<str>datetime_get(e.performed_on.date, 'month'), 2, "0")
  #       )
  #     )),
  #     base := assert_single(e.site.code ++ "_" ++ date_suffix)
  #   select assert_exists(assert_single(base ++ "." ++ <str>number_suffix))
  # );

  type SamplingCodeIndex {
    required site: location::Site;
    required year: int32;
    required month: int32;
    required count: int32;
    constraint exclusive on ((.site, .year, .month));
    index on ((.site, .year, .month));
  }

  type Sampling extending Action {

    trigger index_insert after insert for each do (
      insert SamplingCodeIndex {
        site := __new__.event.site,
        year := <int32>datetime_get(__new__.event.performed_on.date, 'year'),
        month := <int32>datetime_get(__new__.event.performed_on.date, 'month'),
        count := 1
      } unless conflict on ((.site, .year, .month)) else (
        update SamplingCodeIndex set {
          count := .count + 1
        }
      )
    );

    required code : str {
      constraint exclusive;
    };

    # required code : str {
    #   annotation description := "Format : SITE_YEARMONTH.NUMBER. The NUMBER suffix is not appended if the site and month tuple is unique.";
    #   annotation default::example := "SOMESITE_202301 ; SOMESITE_202301.1";

    #   default := "";
    #   constraint exclusive;
    #   rewrite insert, update using (
    #     select (
    #       if (__specified__.code)
    #       then __subject__.code
    #       else sampling_code(__subject__)
    #       # .generated_code ++ "." ++
    #       # <str>(select count(
    #       #   detached Sampling filter .generated_code = __subject__.generated_code
    #       # ))
    #     )
    #   );

    #   # rewrite insert, update using (
    #   #   with collisions := (
    #   #     select count(detached Sampling filter .generated_code = __subject__.generated_code)
    #   #   ),
    #   #   select (
    #   #     if ( collisions > 0)
    #   #     then (select .generated_code ++ "." ++ <str>collisions)
    #   #     else (select .generated_code)
    #   #   )
    #   #   # ),
    #   #   # select (
    #   #   #   if (select exists Sampling filter Sampling.code = generated_code)
    #   #   #   then (
    #   #   #     select generated_code ++ "_" ++
    #   #   #     <str>(select count(Sampling) filter Sampling.code = generated_code)
    #   #   #   )
    #   #   #   else (select generated_code)
    #   #   # )
    #   # );
    # };

    # WARNING : during migration, remove pseudo-field when no sampling was performed
    multi methods: SamplingMethod;
    multi fixatives: default::Conservation; # TODO : what is the relation with conservation of samples down the stream ?

    required sampling_target: SamplingTarget;
    multi target_taxa: taxonomy::Taxon {
      annotation title := "Taxonomic groups that were the target of the sampling effort"
      # handle in application logic
      # constraint expression on (exists .target_taxa) except (.sampling_target != SamplingTarget.Taxa);
    };

    sampling_duration: int32;

    comments: str;

    multi habitats: location::Habitat;
    multi access_points: str;

    multi link samples := .<sampling[is samples::BioMaterial];
    multi link reports := .<sampling[is occurrence::OccurrenceReport];
    multi link external_seqs := .<sampling[is seq::ExternalSequence];

    multi link occurring_taxa := (
      with ext_samples_no_seqs := (select .reports filter not exists .sequences)
      select distinct (
        ext_samples_no_seqs.identification.taxon union
        .external_seqs.identification.taxon union
        .samples.identified_taxa
      )
    );
  }

  # type EventDataset extending default::Auditable {
  #   required name: str {
  #     constraint min_len_value(4);
  #     constraint max_len_value(40);
  #   };

  #   multi samplings: Sampling {
  #     on target delete allow;
  #   };
  #   multi abiotic_measurements: AbioticMeasurement {
  #     on target delete allow;
  #   };
  #   multi spottings: Spotting {
  #     on target delete allow;
  #   };

  #   required multi maintainers: people::Person {
  #     on target delete restrict;
  #   };
  #   multi published_in: reference::Article {
  #     on target delete allow;
  #   };
  #   comments: str;
  # }
}
