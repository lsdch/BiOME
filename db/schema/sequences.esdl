
module seq {

  type BatchRequest extending default::Auditable {
    required label: str;
    required target_gene: Gene;

    requested_on: datetime {
      annotation description := "If empty, the request is a draft and can not be processed yet."
    };
    required requested_by: people::Person;

    multi requested_to: people::Person; # or group of persons ?
    achieved_on: datetime;

    comments: str;
    required multi content: DNA;
  }

  scalar type DNAQuality extending enum<Unknown, Contaminated, Bad, Good>;

  type DNAExtractionMethod extending default::Vocabulary, default::Auditable;

  type DNA extending default::Auditable {
    required chelex_tube: tuple<color:str, number:int16>;
    required specimen: samples::Specimen;
    required code: str { constraint exclusive };
    # maybe autogenerate, talk with Lara
    # maybe use the container as a referential ?

    concentration: float32 {
      annotation title := "DNA concentration in ng/µL";
      constraint min_value(1e-3);
    };

    required multi extracted_by: people::Person;
    required extraction_method: DNAExtractionMethod;
    required extracted_on: tuple<date: datetime, precision: date::DatePrecision> {
      constraint date::required_unless_unknown(.date, .precision);
      rewrite insert, update using (
        (date := date::rewrite_date(.extracted_on), precision :=.extracted_on.precision)
      )
    };

    required quality : DNAQuality;
    required is_empty: bool {
      default := false
    };
    required stored_in: storage::DNAStorage;
    multi link PCRs := .<DNA[is PCR];
    comments: str;
  }


  type Gene extending default::Vocabulary, default::Auditable {
    required motu: bool {
      default := false;
    }; # TODO : discuss with Florian
  }

  abstract type Primer extending default::Auditable {
    required label: str { constraint exclusive };
    required code: str { constraint exclusive };
    sequence: str {
      constraint min_len_value(5);
    };
    description: str;
  }

  type PCRForwardPrimer extending Primer, default::Auditable;
  type PCRReversePrimer extending Primer, default::Auditable;

  type PCRSpecificity extending default::Auditable {
    required label: str { constraint exclusive };
    description: str;
  }

  scalar type PCRQuality extending enum<Failure, Acceptable, Good, Unknown>;

  type PCR extending default::Auditable {
    required DNA: DNA;
    required gene: Gene;
    required code: str { constraint exclusive };
    required forward_primer : PCRForwardPrimer;
    required reverse_primer : PCRReversePrimer;
    required quality: PCRQuality;
    required performed_on: tuple<date: datetime, precision: date::DatePrecision> {
      constraint date::required_unless_unknown(.date, .precision);
      rewrite insert, update using (
        (date := date::rewrite_date(__subject__.performed_on), precision := __subject__.performed_on.precision)
      )
    };
    comments: str;
    multi link chromatograms := .<PCR[is Chromatogram];

    # Questions for Lara :
    # - change quality to be a 0-3 scale, with 0 being a failure, then 1,2 or 3 stars
    # - Nested PCR: do we want the initial PCR to be saved ?
    #     or only the nested one, while keeping the information on the initial primers ?
    # - dilution: should we store the dilution factor as well ?
    # - Utiliser un système de tags combinables pour la spécificité ?
    # - how useful would storing files be ?

    # Notes spécificité :
    # - 3KB: taille de fragment suffisant pour avoir à la fois le 16S et le COI (?)
   }


  type ChromatoPrimer extending Primer, default::Auditable;

  scalar type ChromatoQuality extending enum<Contaminated, Failure, Ok, Unknown>;


  type SequencingInstitute extending default::Auditable {
    required name: str { constraint exclusive };
    comments: str;
  }

  type Chromatogram extending default::Auditable {
    # store chromatogram files ? talk with hosting service about space availability
    required PCR: PCR;
    required YAS_number: str { constraint exclusive };
    required primer: ChromatoPrimer;
    required code: str {
      constraint exclusive;
      rewrite insert using (
        .code ?? (select .YAS_number ++ "|" ++ .primer.code)
      );
    };
    required quality: ChromatoQuality;
    # What is the difference between 'NON DISPONIBLE' and "INCONNU" ?

    required provider: SequencingInstitute;
    comments: str;
    multi link sequences := .<chromatograms[is AssembledSequence]
  }

  abstract type Sequence extending default::Auditable {
    required code: str { constraint exclusive };

    # optional human readable label; use NCBI 'DEFINITION' field when applicable
    label: str;

    sequence: str;

    required gene: Gene;
    comments: str;

    legacy: tuple<id: int32, code: str, alignment_code: str> {
      annotation description := "Legacy identifiers for retrocompatibility with data stored in GOTIT.";
    };

    required category := (
      assert_exists(
        if (__source__ is AssembledSequence) then occurrence::OccurrenceCategory.Internal
        else if (__source__ is ExternalSequence) then occurrence::OccurrenceCategory.External
        else {},
        message := "Occurrence category for seq::Sequence subtype " ++ __source__.__type__.name ++ " is undefined"
      )
    );
  }

  type AssembledSequence extending Sequence {
    required sampling: events::Sampling {
      rewrite insert, update using (
        select .specimen.biomat.sampling
      );
    };

    required identification: occurrence::Identification {
      constraint exclusive;
      on source delete delete target;
    };


    required alignmentCode: str { constraint exclusive };
    accession_number: str {
      annotation description := "The NCBI accession number, if the sequence was uploaded.";
      constraint exclusive;
    };
    # required identification: occurrence::Identification;

    multi published_in: references::Article;
    required multi assembled_by: people::Person;
    required multi chromatograms: Chromatogram;
    required specimen: samples::Specimen;
    required is_reference: bool {
      annotation description := "Whether this sequence should be used as the reference for the identification of the specimen";
    };
    constraint exclusive on ((.specimen, .is_reference)) except (not .is_reference);
  }

  type SeqDB extending default::Auditable, default::Vocabulary {
    link_template: str;
  };

  type SeqReference {
    required db: SeqDB {
      on target delete delete source;
    };
    required accession: str {
      constraint min_len_value(3);
      constraint max_len_value(32);
    };
    required is_origin: bool {
      default := false;
    };

    required code := ( .db.code ++ ":" ++ .accession );

    constraint exclusive on ((.db, .accession));
  };

  scalar type ExtSeqOrigin extending enum<Lab, PersCom, DB>;

  type ExternalSequence extending Sequence, occurrence::Occurrence {

    overloaded required code: str {
      constraint exclusive;
      rewrite update using (
        with suffix := (
          if .origin = ExtSeqOrigin.Lab then "lab"
          else if .origin = ExtSeqOrigin.PersCom then "perscom"
          else (
            with sources := (select .referenced_in filter .is_origin).code
            select array_join(array_agg(sources), "|")
          )
        )
        select .identification.taxon.code ++ "[" ++ .sampling.code ++ "]" ++ .specimen_identifier ++ "|" ++ suffix
      );
    }

    required origin: ExtSeqOrigin;
    source_sample: occurrence::ExternalBioMat;
    published_in: references::Article; # TODO: move this to Occurrence schema

    multi referenced_in: SeqReference {
      constraint exclusive;
    };
    # 🚧 enforce required references when origin = DB in the application layer
    # the following statement is not supported by EdgeDB
    # constraint expression on (.has_references) except (.origin != ExtSeqOrigin.DB);

    required specimen_identifier: str {
      annotation description := "An identifier for the organism from which the sequence was produced, provided in the original source";
    };
    original_taxon: str {
      annotation description := "The verbatim identification provided in the original source.";
    };
  }

  alias SequenceWithType := (
    select Sequence {
      *,
      required sampling := assert_exists(
        [is AssembledSequence].sampling ??
        [is ExternalSequence].sampling
      ),
      required identification := assert_exists(
        [is AssembledSequence].identification ??
        [is ExternalSequence].identification
      ),
    }
  );
}
