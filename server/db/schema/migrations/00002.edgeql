CREATE MIGRATION m1qrh76duzrnlyzkqagufayq6ggj7w5cxhzalg6gefjouabjjwmokq
    ONTO m1ginfxg5lspuxcml7jggba5u2umtojefosbirin53evd5umssvkoq
{
  CREATE GLOBAL default::current_user_id -> std::uuid;
  DROP EXTENSION graphql;
};
