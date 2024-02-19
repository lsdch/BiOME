CREATE MIGRATION m1xzxtbkorylvfsw2vrow4bnaaehnvnq2syat5goftchn5qgqnoxeq
    ONTO m1ue5qhhs7bnv5fut34ghzgbhmqv4ctnpqeflbhdr26owak7w4ftrq
{
  ALTER TYPE people::Person {
      CREATE LINK user := (.<identity[IS people::User]);
  };
};
