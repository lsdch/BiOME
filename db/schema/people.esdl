# Module people holds type definitions related to users, persons and institutions.
#
module people {

  function personByAlias(alias: str) -> Person {
    using (
      select assert_exists(
        Person filter .alias = alias,
        message := "Failed to find person with alias: " ++ alias
      )
    );
  };

  scalar type InstitutionKind extending enum<Lab, FundingAgency, SequencingPlatform, Other>;

  type Institution extending default::Auditable {
    required name: str {
      constraint exclusive;
      constraint min_len_value(3);
      constraint max_len_value(128);
    };
    required code: str {
      constraint exclusive;
      constraint min_len_value(2);
      constraint max_len_value(12);
    };
    description: str {
      rewrite insert, update using (default::null_if_empty(.description));
    };

    required kind: InstitutionKind;

    multi link people := .<institutions[is Person];
  }

  function insert_institution(data: json) -> Institution {
    using (
      insert Institution {
        name := <str>data['name'],
        code := <str>data['code'],
        kind := <InstitutionKind>data['kind'],
        description := <str>json_get(data, 'description'),
      }
    );
  };

  function insert_or_create_institution(data: json) -> Institution {
    using (
      if (json_typeof(data) = 'object')
      then (select (insert_institution(data)))
      else if (json_typeof(data) = 'string') then (
        select (
          assert_exists(
            Institution filter .code = <str>data,
            message := "Failed to find institution with code: " ++ <str>data
          )
        )
      )
      else (
        assert_exists(<Institution>{},
        message := "Invalid institution JSON type: " ++ json_typeof(data))
      )
    );
  };

  type Person extending default::Auditable {
    required first_name: str {
      constraint min_len_value(1);
      constraint max_len_value(30);
    };
    required last_name: str {
      constraint min_len_value(2);
      constraint max_len_value(30);
    };
    required property full_name := array_join([.first_name, .last_name], ' ');
    required alias: str {
      constraint exclusive;
      constraint min_len_value(3);
      constraint max_len_value(32);
      default := <str>{};
      rewrite insert, update using (
        default::null_if_empty(.alias) ?? (
        with
          default_alias := str_lower(.first_name[0] ++ .last_name),
          conflicts := (select count (
            detached Person filter (
              str_trim(.alias, "0123456789") = default_alias
            )
          )),
          suffix := if conflicts > 0 then <str>(conflicts) else "",
        select (default_alias ++ suffix)
        )
      );
    };

    index on ((.alias, .first_name, .last_name));

    contact: str {
      rewrite insert, update using (default::null_if_empty(.contact));
    };
    multi institutions: Institution;

    comment : str;

    link user := .<identity[is User];

    optional property role := .user.role;
  }

  function insert_person(data: json) -> Person {
    using (
      insert Person {
        first_name := <str>data['first_name'],
        last_name := <str>data['last_name'],
        alias := <str>json_get(data, 'alias'),
        contact := <str>json_get(data, 'contact'),
        comment := <str>json_get(data, 'comment'),
        institutions := (
          distinct (for inst in json_array_unpack(json_get(data, 'institutions')) union (
            insert_or_create_institution(inst)
          ))
        )
      }
    );
  };


  scalar type UserRole extending enum<Visitor, Contributor, Maintainer, Admin>;


  type User {
    required login: str {
      constraint exclusive;
      constraint min_len_value(5);
      constraint max_len_value(16);
    };
    required email: str {
      constraint exclusive;
    };

    index on ((.email, .login));

    required password: str {
      annotation description := "Password hashing is done within the database, raw password must be used when creating/updating.";
      rewrite insert, update using (
        if __specified__.password
        then ext::pgcrypto::crypt(.password, ext::pgcrypto::gen_salt('bf', 10))
        else .password
      );
    };

    required role: UserRole {
      default := UserRole.Visitor;
    };

    required identity: Person {
      constraint exclusive;
      on target delete restrict;
      on source delete allow;
    };
  }

  type PendingUserRequest {
    required email: str {
      constraint exclusive;
    };
    index on (.email);
    required first_name: str;
    required last_name: str;
    required full_name := .first_name ++ " " ++ .last_name;
    institution: str;
    motive: str;
    required created_on: datetime {
      rewrite insert using (datetime_of_statement());
    };
    required email_verified: bool {
      default := false;
    };
  }
}