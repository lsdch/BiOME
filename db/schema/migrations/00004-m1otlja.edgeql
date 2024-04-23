CREATE MIGRATION m1otljal5ypsejcqyp4axj4prfz5nw5criwr5e7jop7zxj7r7q5uya
    ONTO m1egbksxywcdlfyuibcvplobjifj2mjdxztlg6rbqvh35ufesmpv3q
{
  ALTER TYPE people::UserInvitation {
      CREATE PROPERTY dest: std::str {
          CREATE REWRITE
              INSERT 
              USING ((IF (NOT (EXISTS (.dest)) OR (.dest = '')) THEN .identity.contact ELSE .dest));
          CREATE REWRITE
              UPDATE 
              USING ((IF (NOT (EXISTS (.dest)) OR (.dest = '')) THEN .identity.contact ELSE .dest));
      };
  };
};
