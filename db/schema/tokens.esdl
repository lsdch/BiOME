module tokens {

  # Tokens are one-time consumables to grant rights on some operations
  abstract type Token {
    required token: str {
      constraint exclusive;
    };
    required expires: datetime;
  }

  type UserInvitation extending Token {
    required identity: people::Person;
    required role: people::UserRole;
    required email: str;
    required issued_by: people::User {
      default := (global default::current_user);
    };
  }

  type PasswordReset extending Token {
    required user: people::User {
      delegated constraint exclusive;
      on target delete delete source;
    };
  };

  type EmailConfirmation extending Token {
    required email: str {
      constraint exclusive;
    };
  };
}