module default {

  # Tracks operations performed on an object
  type Meta {
    annotation title := "Tracking data modifications";

    required created: datetime {
      default := datetime_of_statement()
    };
    created_by_user: people::User {
      default := (select global current_user);
    };
    property created_by := (
      id := .created_by_user.id,
      login := .created_by_user.login,
      name := .created_by_user.identity.full_name,
      alias := .created_by_user.identity.alias
    );


    modified: datetime;
    modified_by_user: people::User {
      rewrite update using (select global current_user)
    };
    property updated_by := (
      id := .modified_by_user.id,
      login := .modified_by_user.login,
      name := .modified_by_user.identity.full_name,
      alias := .modified_by_user.identity.alias
    );

    property lastUpdated := (.modified ?? .created);
  }

  # Operations on auditable objects are automatically tracked
  abstract type Auditable {
    annotation title := "Auto-generation of timestamps";

    required meta: Meta {
      constraint exclusive;
      on source delete delete target;
      on target delete restrict;
      default := (insert Meta {});
      rewrite update using (
        update .meta set { modified := datetime_of_statement() }
      );
    }
  }
}