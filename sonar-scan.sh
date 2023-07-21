 #!/bin/bash
 
 ~/go/bin/sonar-scanner \
  -Dsonar.projectKey=numbered-notation-xml \
  -Dsonar.sources=. \
  -Dsonar.host.url=http://localhost:9000 \
  -Dsonar.token=sqp_5f6495f1d551a25142e10682374ee08d0ed6252d