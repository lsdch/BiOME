CREATE MIGRATION m1gjaelgh4qeg5yu3g63zaisym6ba57er3lro5gzkuhilgarzx3coq
    ONTO m1vz4wlfzuhbpvm64nxquge37sq7hdiiy4yb3fpaodsdqf4detf2za
{
  ALTER TYPE tokens::EmailConfirmation RENAME TO tokens::EmailVerification;
};
