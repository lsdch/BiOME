CREATE MIGRATION m1gwdc3fyv5f5evagrf6nmvjjuxy54mwqygepqyzegij5ob2idxrbq
    ONTO m1apy5rrwgvslzo7myk37d63mwsidhdllsw55hxtpmebdrul6rzrsq
{
  ALTER TYPE event::Sampling {
      ALTER LINK occurences {
          RENAME TO occurrences;
      };
  };
};
