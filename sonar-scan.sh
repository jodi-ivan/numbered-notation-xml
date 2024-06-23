 #!/bin/bash
 
 ~/go/bin/sonar-scanner \
  -Dsonar.projectKey=numbered-notation-xml \
  -Dsonar.sources=. \
  -Dsonar.host.url=http://localhost:9000 \
  -Dsonar.token=sqp_bcc0763c5f0eb77048397da13de7668a84e99d9d