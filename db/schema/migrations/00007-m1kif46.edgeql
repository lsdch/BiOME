CREATE MIGRATION m1kif46mbpxfhrweaq7aws2fueqshth6wopscxrwznuvezgslcju4q
    ONTO m1dnczj4wfucbo4z7bq4s3fnwijc3ywkeh6xgyvhdbemnbgm65f7oa
{
  DROP ALIAS datasets::AnyDataset;
  ALTER TYPE datasets::AbstractDataset RENAME TO datasets::Dataset;
  ALTER TYPE datasets::Dataset {
      CREATE REQUIRED PROPERTY category := (<datasets::DatasetCategory>std::str_replace((std::str_split(.__type__.name, '::'))[1], 'Dataset', ''));
  };
};
