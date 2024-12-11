module references {

  function generate_article_code(authors: array<str>, year: int32) -> str using (
    select (
      if len(authors) = 1
      then str_split(authors[0], ' ')[0] ++ "_" ++ <str>year
      else if len(authors) = 2
      then (
        array_join([
          str_split(authors[0], ' ')[0],
          str_split(authors[1], ' ')[0],
        ]
        , "_") ++ "_" ++ <str>year
      )
      else str_split(authors[0], ' ')[0] ++ "_et_al_" ++ <str>year
    )
  );

  global alphabet := "abcdefghijklmnopqrstuvxyz";

  type Article extending default::Auditable {
    required authors: array<str>;
    required year: int32 {
      constraint min_value(1000);
    };
    title: str;
    journal: str;
    verbatim: str;
    doi: str;
    comments: str;
    required code: str {
      constraint exclusive;
      default := (
        with
          gen_code := generate_article_code(.authors, .year),
          discriminant := (select count(detached Article filter generate_article_code(.authors, .year) = gen_code) - 1),
        select (
          if discriminant >= 0 then generate_article_code(.authors, .year) ++ global alphabet[discriminant]
          else generate_article_code(.authors, .year)
        )
      );
    };

    index on ((.code, .year));
  }
}