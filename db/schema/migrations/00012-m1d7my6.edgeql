CREATE MIGRATION m1d7my64slzvskmfk5zhhdjinqpepbhshkjplb6flbs4uxqkfmrd4q
    ONTO m1wytcfnhahhtbe72lczcaaueqjdyazr7b65x2vecflwvpvb5qqcoq
{
  CREATE MODULE tokens IF NOT EXISTS;
  ALTER TYPE people::Token RENAME TO tokens::Token;
  ALTER TYPE people::EmailConfirmation RENAME TO tokens::EmailConfirmation;
  ALTER TYPE people::PasswordReset RENAME TO tokens::PasswordReset;
  ALTER TYPE people::UserInvitation RENAME TO tokens::UserInvitation;
};
