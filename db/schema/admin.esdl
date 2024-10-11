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

    required instance := assert_exists((select InstanceSettings limit 1));

    email := (select EmailSettings limit 1);

    required security:= assert_exists((select SecuritySettings limit 1));
  }

  type SecuritySettings {
    required min_password_strength: int32 {
      default := 3;
      constraint min_value(3);
      constraint max_value(5);
    };
    required refresh_token_lifetime: int32 {
      annotation description := "Validity period for a session refresh token in hours";
      constraint min_value(1);
      default := 1;
    };
    required invitation_token_lifetime: int32 {
      annotation description := "Validity period for an account invitation in days";
      constraint min_value(1);
      default := 7;
    }
    required jwt_secret_key: str {
      constraint min_len_value(32);
    };
  }
}