module datasets {

  type Alignment extending default::Auditable {
    required label: str;
    required multi sequences: seq::Sequence;
    comments: str;
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
    required number: int16;
    required method: DelimitationMethod;
    required multi sequences: seq::Sequence;
  }
}