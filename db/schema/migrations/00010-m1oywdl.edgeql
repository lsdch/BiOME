CREATE MIGRATION m1oywdlm6aub3uhvikw4qu34c2bblgfjxcw4tazbuckglznpay5irq
    ONTO m17zsgnd5mxympn77nxxeighgwo4hljbdarc6s6udbrmhextikmakq
{
  ALTER FUNCTION taxonomy::taxon_code(taxon: taxonomy::Taxon) USING (std::re_replace('[^()a-zA-Z0-9]+', '_', std::str_replace(taxon.name, 'sp. ', 'sp.'), flags := 'ig'));
};
