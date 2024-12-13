module events {

  type Program extending default::Auditable {
    required label: str { constraint exclusive };
    required code: str { constraint exclusive };
    index on ((.code, .label));

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

    required code := (
      with
        date := .performed_on.date,
        precision := .performed_on.precision
      select .site.code ++ "|" ++ (
        if precision = date::DatePrecision.Unknown then "undated"
        else if precision = date::DatePrecision.Year then <str>datetime_get(date, 'year')
        else <str>datetime_get(date, 'year') ++ "-" ++ <str>datetime_get(date, 'month')
      )
    );

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
    required param: AbioticParameter {
      on target delete delete source;
    };
    required value: float32;
    constraint exclusive on ((.event, .param))
  }

  scalar type SamplingTarget extending enum<Community, Unknown, Taxa>;

  type SamplingMethod extending default::Vocabulary, default::Auditable;

  scalar type SamplingNumber extending sequence;

  type Sampling extending Action {

    required number: SamplingNumber {
      readonly := true;
    };

    required single code := (
      with
        id := .id,
        event := .event,
        sisters := (
          select detached events::Sampling
          filter .event.code = event.code
          order by .number
        ),
        rank := (select assert_single(enumerate(sisters) filter .1.id = id).0),
        suffix := (if count(sisters) > 1 then '.' ++ <str>(rank + 1) else '')
        select assert_single(assert_exists(event.code ++ suffix))
    );

    # WARNING : during migration, remove pseudo-field when no sampling was performed
    multi methods: SamplingMethod;
    multi fixatives: samples::Fixative; # TODO : what is the relation with conservation of samples down the stream ?

    required sampling_target: SamplingTarget;
    multi target_taxa: taxonomy::Taxon {
      annotation title := "Taxonomic groups that were the target of the sampling effort"
      # handle in application logic
      # constraint expression on (exists .target_taxa) except (.sampling_target != SamplingTarget.Taxa);
    };

    # sampling duration in minutes
    sampling_duration: int32;

    comments: str;

    multi habitats: sampling::Habitat;
    multi access_points: str;

    multi link samples := .<sampling[is occurrence::BioMaterial];
    multi link external_seqs := .<sampling[is seq::ExternalSequence];

    multi link occurring_taxa := (
      with ext_samples_no_seqs := (
        select .samples[is occurrence::ExternalBioMat]
        filter not exists [is occurrence::ExternalBioMat].sequences
      )
      select distinct (
        ext_samples_no_seqs.identification.taxon union
        .external_seqs.identification.taxon union
        .samples[is occurrence::InternalBioMat].identified_taxa
      )
    );
  }
}
