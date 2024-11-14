module samples {

  type BioMaterial extending occurrence::Occurrence {
    required code : str {
      constraint exclusive;
      annotation description := "Format like 'taxon_code|sampling_code'";
      rewrite insert, update using (
        .identification.taxon.code ++ "|" ++ events::event_code(.sampling.event)
      )
    };

    required created_on: tuple<date: datetime, precision: date::DatePrecision> {
      constraint date::required_unless_unknown(.date, .precision);
      rewrite insert, update using (
        (date := date::rewrite_date(.created_on), precision :=.created_on.precision)
      )
    };
    required multi sorted_by: people::Person;
    multi published_in: references::Article;
    multi link content := .<biomat[is Sample];
    multi link specimens := .<biomat[is Specimen];
    multi link bundles := .<biomat[is BundledSpecimens];
    multi link identified_taxa := (select distinct .specimens.identification.taxon ?? .identification.taxon);
  }

  type ContentType extending default::Vocabulary, default::Auditable;

  abstract type Sample extending default::Auditable {
    required biomat: BioMaterial;
    required type: ContentType;
    required conservation: default::Conservation;

    required number: int16 {
      annotation description := "Incremental number that discriminates between tubes having the same type in a bio material lot. Used to generate the tube code.";
      rewrite insert using (
        select 1 + (
          max(Sample.number
            filter Sample.biomat = __subject__.biomat
            and Sample.type = __subject__.type
          ) ?? 0)
      )
    };
    property tube := (.type.code ++ <str>.number);
    # required tube: tuple<number: int16, code: str> {
    #   rewrite insert using (
    #     with biomat := __subject__.biomat,
    #     type := __subject__.type,
    #     number := (select 1 + (
    #       max(Sample.tube.number
    #         filter Sample.biomat = biomat
    #         and Sample.type = type
    #       ) ?? 0))
    #     select (number := number, code := .type.code ++ <str>number)
    #   );
    # }
    comments: str;
  }

  type BundledSpecimens extending Sample {
    annotation description := "A tube containing several specimens.";
    required quantity: int16 {
      constraint min_value(2)
    };
  }

  type Specimen extending Sample {
    annotation description := "A single specimen isolated in a tube.";
    required morphological_code: str {
      constraint exclusive;
      rewrite insert, update using (
        (__subject__.biomat.code ++ "[" ++ .tube ++ "]")
      )
    };

    dissected_by: people::Person;

    molecular_code: str {
      constraint exclusive;
    }; # auto-generated
    required molecular_number: str {
      rewrite insert using (
        <str>(select count(Specimen) filter Specimen.biomat.sampling.event.site = __subject__.biomat.sampling.event.site)
      );
      # is string because retrocompatibility
    };

    multi link sequences := .<specimen[is seq::AssembledSequence];
    link identification := (
      (select .sequences filter .is_reference).identification
    );

    multi pictures: default::Picture {
      on source delete delete target;
      # APP: remember to also delete files
    };
    multi link slides := .<specimen[is Slide];

  }

  type Slide extending default::Auditable {
    required specimen: Specimen;
    required label: str; # the physical label
    property code: str {
      annotation description := "Generated as '{collectionCode}_{containerCode}_{slidePositionInBox}'";
      constraint exclusive;
      rewrite insert, update using (
        (select array_join([.storage.collection.code, .storage.code, <str>.storage_position], "_"))
      )
    };

    required created_on: tuple<date: datetime, precision: date::DatePrecision> {
      constraint date::required_unless_unknown(.date, .precision);
      rewrite insert, update using (
        (date := date::rewrite_date(.created_on), precision :=.created_on.precision)
      )
    };

    required multi mounted_by: people::Person;

    required storage: storage::SlideStorage;
    required storage_position: int16;
    constraint exclusive on ((.storage, .storage_position));

    multi pictures: default::Picture {
      on source delete delete target;
    };
    comment: str;
  }
}
