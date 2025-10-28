CREATE MIGRATION m1ha3vu5s6nlivik2ecyyaomjpewtqlfvo5e575uxzp7fqurupszma
    ONTO m1w5obxmcerof4vfqdom2kqv5vp6zs27ac26hzuhztlcejwz652epa
{
  ALTER TYPE datasets::OccurrenceDataset {
      DROP PROPERTY is_congruent;
  };
};
