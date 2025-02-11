CREATE MIGRATION m1k6kg5dfxgv4xdyzim2gfyhyyebrv4ab7ep3e3y7qb4j2gbfme5za
    ONTO m1uicsjoj2sfqrubiqwsyhoxyfsl6wh2ggr7flrhfyrlncpo6seyga
{
  CREATE FUNCTION default::get_vocabulary(code: std::str, object_name: OPTIONAL std::str = 'vocabulary') ->  default::Vocabulary USING (std::assert_exists(std::assert_single((SELECT
      default::Vocabulary
  FILTER
      (.code = code)
  )), message := ((('Failed to find ' ++ object_name) ++ ' with code: ') ++ code)));
  ALTER TYPE occurrence::ExternalBioMat {
      CREATE LINK original_source: references::DataSource;
  };
};
