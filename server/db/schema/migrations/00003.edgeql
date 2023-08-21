CREATE MIGRATION m1pjrvyeddw466xedtztu3rfjux6vjxy5q3rdgjvgfei6gbeor2gba
    ONTO m1nq6igapvk5az34hdeuik4vgqwy6kfgra73qleshkrzrmedjfkzlq
{
  ALTER TYPE event::PlannedSampling {
      ALTER LINK targetTaxa {
          RENAME TO target_taxa;
      };
  };
  ALTER TYPE event::Sampling {
      ALTER LINK biomaterials {
          RENAME TO samples;
      };
  };
  ALTER TYPE event::Sampling {
      CREATE MULTI LINK occurences := (.<sampling[IS occurrence::Occurrence]);
      CREATE MULTI LINK reports := (.<sampling[IS occurrence::OccurrenceReport]);
  };
  ALTER TYPE occurrence::Dataset {
      ALTER LINK events {
          RENAME TO samplings;
      };
  };
  ALTER TYPE samples::Slide {
      ALTER PROPERTY code {
          CREATE ANNOTATION std::description := "Generated as '{collectionCode}_{containerCode}_{slidePositionInBox}'";
      };
  };
};
