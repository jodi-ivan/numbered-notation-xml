 #!/bin/bash
 
 ~/go/bin/sonar-scanner \
  -Dsonar.projectKey=numbered-notation-xml \
  -Dsonar.sources=. \
  -Dsonar.host.url=http://localhost:9000 \
  -Dsonar.token=sqp_3cd177b0e6e6d9e9a5918783f33561ec5e1f2012