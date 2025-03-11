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

    index on (.email);
  }

  type PasswordReset extending Token {
    required user: people::User {
      delegated constraint exclusive;
      on target delete delete source;
    };
  };

  type EmailVerification extending Token {
    required user_request: people::PendingUserRequest {
      constraint exclusive;
      on target delete delete source;
    };
  };

  type SessionRefreshToken extending Token {
    required user: people::User {
      on target delete delete source;
    };
  }
}