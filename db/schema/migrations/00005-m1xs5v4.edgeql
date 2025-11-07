CREATE MIGRATION m1xs5v4zu7mg5woszun3ds53ti7scatdkdpgja54ea47zx5u5on34q
    ONTO m1idi3q4khw76q5d3nw46a43gghhizszyoeqn2x7ylhyxrubxyyuga
{
  ALTER TYPE events::Sampling {
      DROP PROPERTY sampling_target;
  };
  DROP SCALAR TYPE events::SamplingTarget;
};
