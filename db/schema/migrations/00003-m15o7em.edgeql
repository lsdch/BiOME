CREATE MIGRATION m15o7em7rbkbulecfov4lp2j3ru5ncgjz3sgcji77246tglbbc26gq
    ONTO m1c3enlyoxfvhjbyemzecykb6h77ofsfihvmc7duqtxwqgtbplh6aa
{
  ALTER TYPE occurrence::Occurrence {
      ALTER PROPERTY category {
          USING (.__type__.name);
      };
  };
  ALTER TYPE seq::Sequence {
      ALTER PROPERTY category {
          USING (.__type__.name);
      };
  };
};
