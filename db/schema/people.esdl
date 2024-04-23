# Module people holds type definitions related to users, persons and institutions.
#
module people {

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

    index on (.code);
  }

  type Person extending default::Auditable {
    required first_name: str {
      constraint min_len_value(1);
      constraint max_len_value(30);
    };
    required last_name: str {
      constraint min_len_value(2);
      constraint max_len_value(30);
    };
    required property full_name := array_join(
      array_agg({.first_name, .last_name}), ' '
    );
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
    }

    contact: str {
      rewrite insert, update using (default::null_if_empty(.contact));
    };
    multi institutions: Institution;

    comment : str;

    link user := .<identity[is User];

    optional property role := .user.role;
  }

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

    required password: str {
      annotation description := "Password hashing is done within the database, raw password must be used when creating/updating.";
      rewrite insert, update using (
        if __specified__.password
        then ext::pgcrypto::crypt(.password, ext::pgcrypto::gen_salt('des'))
        else .password
      );
    };

    required email_confirmed: bool {
      default := false
    };

    required role: UserRole {
      default := UserRole.Visitor;
    };

    identity: Person {
      constraint exclusive;
      on target delete restrict;
      on source delete allow;
    };

    is_active := (.email_confirmed and exists .identity);
  }

  alias ActiveUser := (select User filter .is_active);
  alias InactiveUser := (select User filter not .is_active);

  type PendingUserRequest {
    required user: User {
      constraint exclusive;
      on target delete delete source;
      on source delete allow;
    };
    required identity: tuple<
      first_name: str,
      last_name: str,
      institution: str
    >;
    required motive: str;
    required created_on: datetime {
      rewrite insert using (datetime_of_statement());
    };
  }

  type UserInvitation {
    required identity: Person;
    required role: UserRole;
    required dest: str;
    required token: str {
      constraint exclusive;
    };
    required expires: datetime;
    required issued_by: User {
      default := (global default::current_user);
    };
  }

  # AccountToken is a consumable token that confirms that some operation
  # on an account is allowed.
  abstract type AccountToken {
    required token: str {
      constraint exclusive;
    };
    required expires: datetime;
    required user: User {
      delegated constraint exclusive;
      on target delete delete source;
    };
  }

  type PasswordReset extending AccountToken;

  type EmailConfirmation extending AccountToken;
}