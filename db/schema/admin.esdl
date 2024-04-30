module admin {

  type InstanceSettings {
    required name: str {
      default := "DarCo";
      constraint min_len_value(3);
      constraint max_len_value(20);
    }

    description: str;

    required public: bool {
      annotation description := "Whether parts of the platform are open to the public (anonymous users).";
      default := true;
    };

    required allow_contributor_signup: bool {
      annotation description := "Whether a user can request a new account with contributor privileges";
      default := true;
    };
  }

  type EmailSettings {
    required host: str;
    required user: str;
    required password: str;
    required port: int32;
  }

  type Settings {

    required instance: InstanceSettings {
      constraint exclusive;
      on source delete delete target;
    };


    email: EmailSettings {
      on source delete delete target;
    };

    required security: SecuritySettings {
      constraint exclusive;
      on source delete delete target;
    };
  }

  type SecuritySettings {
    required min_password_strength: int32 {
      default := 3;
      constraint min_value(3);
      constraint max_value(5);
    };
    required auth_token_lifetime: int32 {
      annotation description := "Validity period for an authentication token in minutes";
      constraint min_value(5);
      default := 15;
    };
    required account_token_lifetime: int32 {
      annotation description := "Validity period for an account operation token in hours";
      constraint min_value(1);
      default := 24;
    };
    required jwt_secret_key: str {
      constraint min_len_value(32);
    };
  }
}