#!/usr/bin/env bash
make && cat res/json/ncbi-SAMN03894263.json | ./bscurate -config=res/config.json -schema=res/schemas/config-schema.json
curl https://www.ebi.ac.uk/biosamples/samples/SAMN03894263.json -o in.json && cat in.json | ./bscurate -config=res/config.json -schema=res/schemas/config-schema.json
curl https://www.ebi.ac.uk/biosamples/samples/SAMEA3774859.json -o in.json && cat in.json | ./bscurate -config=res/config.json -schema=res/schemas/config-schema.json
curl https://wwwdev.ebi.ac.uk/biosamples/samples/SAMEA4669076.json -o in.json && cat in.json | ./bscurate -config=res/config.json -schema=res/schemas/config-schema.json