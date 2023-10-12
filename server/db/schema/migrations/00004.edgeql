CREATE MIGRATION m16dj7xaefgahqazq34efjxkudsj3mpmhk5rdmolnrwlk3hgbov2gq
    ONTO m1cq5kdmkhjv26n2ndnguc36du5vv4kyvqq35jok7qngla2vdmv45a
{
  ALTER TYPE event::AbioticParameter {
      CREATE REQUIRED PROPERTY unit: std::str {
          SET REQUIRED USING (<std::str>{});
      };
  };
};
