CREATE MIGRATION m1dnczj4wfucbo4z7bq4s3fnwijc3ywkeh6xgyvhdbemnbgm65f7oa
    ONTO m1o6pb7ep6qf7qn5irwkmuowhyeendjrs3c2xdoxp36nxozazkqu4q
{
  ALTER TYPE datasets::AbstractDataset {
      CREATE REQUIRED PROPERTY pinned: std::bool {
          SET default := false;
      };
  };
};
