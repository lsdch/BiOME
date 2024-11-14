module references {
  type Article extending default::Auditable {
    multi authors: str;
    required year: int32 {
      constraint min_value(1000);
    };
    required title: str;
    journal: str;
    verbatim: str;
    doi: str;
    comments: str;
   }
}