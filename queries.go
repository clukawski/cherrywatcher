package main

// Not dead yet
const cherryRIPQuery = `https://query.wikidata.org/sparql?query=SELECT%20distinct%20%3Fitem%20%3FitemLabel%20(SAMPLE(%3FRIP)%20as%20%3FRIP)%20WHERE%20%7B%0A%20%20%3Fitem%20wdt%3AP31%20wd%3AQ5.%0A%20%20%3Fitem%20%3Flabel%20%22Don%20Cherry%22%40en.%0A%20%20OPTIONAL%7B%3Fitem%20wdt%3AP570%20%3FRIP%20.%7D%20%20%20%20%20%23%20P570%20%3A%20Date%20of%20death%0A%0A%20%20SERVICE%20wikibase%3Alabel%20%7B%20bd%3AserviceParam%20wikibase%3Alanguage%20%22en%22.%20%7D%0A%20%20FILTER%20(%20%3Fitem%20in%20(%20wd%3AQ592524%20)%20)%0A%7D%0AGROUP%20BY%20%3Fitem%20%3FitemLabel%20%3FitemDescription`

// This is some other Don Cherry, but on the upside they're dead so we can test this code
const otherCherryRIPQuery = `https://query.wikidata.org/sparql?query=SELECT%20distinct%20%3Fitem%20%3FitemLabel%20(SAMPLE(%3FRIP)%20as%20%3FRIP)%20WHERE%20%7B%0A%20%20%3Fitem%20wdt%3AP31%20wd%3AQ5.%0A%20%20%3Fitem%20%3Flabel%20%22Don%20Cherry%22%40en.%0A%20%20OPTIONAL%7B%3Fitem%20wdt%3AP570%20%3FRIP%20.%7D%20%20%20%20%20%23%20P570%20%3A%20Date%20of%20death%0A%0A%20%20SERVICE%20wikibase%3Alabel%20%7B%20bd%3AserviceParam%20wikibase%3Alanguage%20%22en%22.%20%7D%0A%23%20%20FILTER%20(%20%3Fitem%20in%20(%20wd%3AQ592524%20)%20)%0A%20%20FILTER%20(%20%3Fitem%20in%20(%20wd%3AQ456180%20)%20)%0A%7D%0AGROUP%20BY%20%3Fitem%20%3FitemLabel%20%3FitemDescription`

// Underlying query:
//
// SELECT distinct ?item ?itemLabel (SAMPLE(?RIP) as ?RIP) WHERE {
//   ?item wdt:P31 wd:Q5.
//   ?item ?label "Don Cherry"@en.
//   OPTIONAL{?item wdt:P570 ?RIP.} # P570 : Date of death
//   SERVICE wikibase:label { bd:serviceParam wikibase:language "en". }
//   FILTER ( ?item in ( wd:Q592524 ) ) # wd:Q592524 : Don Cherry entity ID, filter out other Don Cherries
// }
// GROUP BY ?item ?itemLabel ?itemDescription
