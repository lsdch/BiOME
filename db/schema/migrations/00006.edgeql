CREATE MIGRATION m1ue5qhhs7bnv5fut34ghzgbhmqv4ctnpqeflbhdr26owak7w4ftrq
    ONTO m1cop5taizv5aeqn6t4uzzhbicko52fssrg25ddjhrvv5kifbf5jba
{
  ALTER TYPE people::Institution {
      CREATE INDEX ON (.code);
  };
};
