using extension pgcrypto;
using extension auth;
using extension pg_trgm;



module default {

  global current_user_id: uuid {
    default := <uuid>{};
  };

  # Dynamically computed upon querying
  # Is an empty set if user is not authenticated
  global current_user := (
    select people::User filter .id = global current_user_id
  );

  abstract annotation example;

  # Trims a string, and transforms to an empty set if result is the null string
  function null_if_empty(s: str) -> optional str using(
    with trimmed := str_trim(s)
    select <str>{} if len(trimmed) = 0 else trimmed
  );

  abstract type Vocabulary {
    annotation title := "An extensible list of terms";

    required label: str {
      constraint exclusive;
      constraint min_len_value(4);
    };
    required code: str {
      annotation title := "An expressive, unique, user-generated uppercase alphanumeric code";
      delegated constraint exclusive;
      constraint min_len_value(2);
      constraint max_len_value(12);
    };
    description: str;
  }

  type Picture {
    legend: str;
    required path: str {
      constraint exclusive;
    };
  }
}



