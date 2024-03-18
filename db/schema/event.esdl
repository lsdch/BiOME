module event {

  type Program extending default::Auditable {
    required label: str { constraint exclusive };
    required code: str { constraint exclusive };
    required multi managers: people::Person;
    multi funding_agencies: people::Institution;
    start_year: int16 {
      constraint min_value(1900);
    };
    end_year: int16;
    constraint expression on (.start_year <= .end_year);
    comments: str;
  }

  abstract type Event extending default::Auditable {
    required site: location::Site;

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
  }

  type Spotting extending Event {
    multi target_taxa: taxonomy::Taxon;
  }

  type AbioticParameter extending default::Vocabulary, default::Auditable {
    required unit: str;
  };

  type AbioticMeasurement extending event::Event {
    required param: AbioticParameter;
    required value: float32;
    related_sampling: Sampling;
    constraint exclusive on ((.site, .param, .performed_on))
  }

  scalar type SamplingTarget extending enum<Community, Unknown, Taxa>;

  type SamplingMethod extending default::Vocabulary, default::Auditable;

  type Sampling extending Event {
    property generated_code := (
      select .site.code ++
      "_" ++ <str>datetime_get(.performed_on.date, 'year') ++
      <str>datetime_get(.performed_on.date, 'month')
    );
    required code : str {
      annotation title := "Unique sampling identifier, auto-generated at sampling creation.";
      annotation description := "Format : SITE_YEARMONTH_NUMBER. The NUMBER suffix is not appended if the site and month tuple is unique.";
      annotation default::example := "SOMESITE_202301 ; SOMESITE_202301_1";

      constraint exclusive;
      rewrite insert using (
        (select
          .generated_code ++ "_" ++
          <str>(select count(Sampling) filter Sampling.code = __subject__.generated_code
        ) if (select exists Sampling filter Sampling.code = __subject__.generated_code)
          else .generated_code)
      )
    };

    # WARNING : during migration, remove pseudo-field when no sampling was performed
    multi methods: SamplingMethod;
    multi fixatives: default::Conservation; # TODO : what is the relation with conservation of samples down the stream ?

    required sampling_target: SamplingTarget;
    multi target_taxa: taxonomy::Taxon {
      annotation title := "Taxonomic groups that were the target of the sampling effort"
      # handle in application logic
      # constraint expression on (exists .target_taxa) except (.sampling_target != SamplingTarget.Taxa);
    };

    sampling_duration: duration;

    required is_donation: bool;

    comments: str;

    multi link measurements := .<related_sampling[is AbioticMeasurement];
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

  type EventDataset extending default::Auditable {
    required name: str {
      constraint min_len_value(4);
      constraint max_len_value(40);
    };

    multi samplings: Sampling {
      on target delete allow;
    };
    multi abiotic_measurements: AbioticMeasurement {
      on target delete allow;
    };
    multi spottings: Spotting {
      on target delete allow;
    };

    required multi maintainers: people::Person {
      on target delete restrict;
    };
    multi published_in: reference::Article {
      on target delete allow;
    };
    comments: str;
  }
}
