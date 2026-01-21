CREATE MIGRATION m17zsgnd5mxympn77nxxeighgwo4hljbdarc6s6udbrmhextikmakq
    ONTO m1zojeutj47b5airchz4mwkcrzggkwmb2q3oat6vesmuhvyt4f27za
{
  ALTER FUNCTION taxonomy::taxon_code(taxon: taxonomy::Taxon) USING (std::re_replace('[^()a-zA-Z0-9]', '_', std::str_replace(taxon.name, 'sp. ', 'sp.'), flags := 'ig'));
};
