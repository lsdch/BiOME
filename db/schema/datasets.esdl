module datasets {

  type AbstractDataset extending default::Auditable {
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

  # type SiteDataset extending AbstractDataset {
  #   multi link sites := .<datasets[is location::Site];
  # }

  type Dataset extending AbstractDataset {
    multi link sites := .<datasets[is location::Site];
  }

  type SeqDataset extending AbstractDataset {
    required multi sequences: seq::Sequence;
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