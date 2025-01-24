module datasets {

  abstract type AbstractDataset extending default::Auditable {
    required label: str {
      constraint min_len_value(4);
      constraint max_len_value(40);
    };
    required slug: str {
      constraint exclusive;
    };
    description: str;

    required multi maintainers: people::Person {
      # edgedb error: SchemaDefinitionError:
      # cannot specify a rewrite for link 'maintainers' of object type 'location::SiteDataset' because it is multi
      # Hint: this is a temporary implementation restriction
      # Currently handling this in the application layer
      # rewrite insert using (global default::current_user.identity union .maintainers);
      # rewrite update using (.maintainers union .meta.created_by_user.identity);
    };
  }

  type SiteDataset extending AbstractDataset {
    multi sites: location::Site {
      on target delete allow;
      on source delete allow;
    };
  };

  type OccurrenceDataset extending AbstractDataset {
    multi occurrences: occurrence::Occurrence {
      on target delete allow;
      on source delete allow;
    };
    multi sites := (.occurrences.sampling.event.site);

    # Dataset is congruent if all occurrences are congruent:
    # Morphological identification matches molecular identification
    required is_congruent := (all(
      (select .occurrences {
        is_congruent := (
          [is occurrence::ExternalBioMat].is_congruent ??
          [is occurrence::InternalBioMat].is_congruent ??
          true
        )
      }).is_congruent
    ));
  }

  type SeqDataset extending AbstractDataset {
    multi sequences: seq::Sequence {
      on target delete allow;
      on source delete allow;
    };
    multi sites := (.sequences.sampling.event.site)
  }


  type MOTUDataset extending default::Auditable {
    required label: str;
    multi link MOTUs := .<dataset[is MOTU];
    published_in: references::Article;
    required multi generated_by: people::Person;
    comments: str;
  }

  type DelimitationMethod extending default::Vocabulary, default::Auditable;

  type MOTU extending default::Auditable {
    required dataset: MOTUDataset;
    required number: int32;
    required method: DelimitationMethod;
    required multi sequences: seq::Sequence;
  }
}