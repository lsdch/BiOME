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

  function event_code(e: Event) -> str using (
    with
      date_suffix := (select assert_single(
        if(e.performed_on.precision = date::DatePrecision.Unknown)
        then "UNK"
        else if (e.performed_on.precision = date::DatePrecision.Year)
        then <str>datetime_get(e.performed_on.date, 'year')
        else (
          <str>datetime_get(e.performed_on.date, 'year') ++
          str_pad_start(<str>datetime_get(e.performed_on.date, 'month'), 2, "0")
        )
      )),
    select (e.site.code ++ "_" ++ date_suffix)
  );

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

    dataset: datasets::Dataset;

    spotting := .<event[is Spotting];
    multi abiotic_measurements := .<event[is AbioticMeasurement];
    multi samplings := .<event[is Sampling];

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

  type Sampling extending Action {

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

    multi habitats: sampling::Habitat;
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
}
