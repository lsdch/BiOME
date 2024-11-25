CREATE MIGRATION m1wyz3xy7mpbhpvqcg56pwsxy76grsz5ttjgcjdbpxl3dap5hb6ifq
    ONTO m1ok2r24epmn57wmjwloel7wqsr3dqtxcvmcvsyodnatqllruzn4wq
{
  ALTER TYPE occurrence::Identification {
      CREATE REQUIRED PROPERTY is_type: std::bool {
          SET default := false;
      };
  };
  ALTER TYPE seq::Sequence {
      CREATE PROPERTY label: std::str;
      CREATE PROPERTY sequence: std::str;
  };
  ALTER TYPE seq::ExternalSequence {
      CREATE LINK reference: references::Article;
  };
};
