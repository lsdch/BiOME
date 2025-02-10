CREATE MIGRATION m1o6pb7ep6qf7qn5irwkmuowhyeendjrs3c2xdoxp36nxozazkqu4q
    ONTO m1qm23x4gbpwnqjw72h3htgve37llf7utfxqxllb44dimjzzulk2tq
{
  CREATE SCALAR TYPE datasets::DatasetCategory EXTENDING enum<Site, Occurrence, Seq>;
  CREATE ALIAS datasets::AnyDataset := (
      SELECT
          datasets::AbstractDataset {
              **,
              category := std::assert_exists((IF (datasets::AbstractDataset IS datasets::SiteDataset) THEN datasets::DatasetCategory.Site ELSE (IF (datasets::AbstractDataset IS datasets::OccurrenceDataset) THEN datasets::DatasetCategory.Occurrence ELSE (IF (datasets::AbstractDataset IS datasets::SeqDataset) THEN datasets::DatasetCategory.Seq ELSE <datasets::DatasetCategory>{}))))
          }
  );
};
