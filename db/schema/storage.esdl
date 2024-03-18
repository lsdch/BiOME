
module storage {

  abstract type Storage extending default::Auditable {
    required label: str {
      constraint exclusive;
      constraint min_len_value(4);
    };
    required code: str {
      constraint exclusive
    }
    description: str;
    required collection: Collection;
  }

  type BioMatStorage extending Storage;
  type SlideStorage extending Storage;
  type DNAStorage extending Storage;

  type Collection {
    required label: str {
      constraint exclusive;
      constraint min_len_value(4);
    };
    required code: str {
      constraint exclusive;
      constraint min_len_value(4)
    };
    required taxon: taxonomy::Taxon;
    required maintainers: people::Person;
    comments: str;

    multi link bio_mat_storages := .<collection[is BioMatStorage];
    multi link slide_storages := .<collection[is SlideStorage];
    multi link DNA_storages := .<collection[is DNAStorage];
  }
}
