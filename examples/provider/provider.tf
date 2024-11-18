## Example: Authenticate with username and password

terraform {
  required_providers {
    sonarqube = {
      source = "jdamata/sonarqube"
    }
  }
}

provider "sonarqube" {
  user = "admin"
  pass = "admin"
  host = "http://127.0.0.1:9000"
}

## Example: Authenticate with token

terraform {
  required_providers {
    sonarqube = {
      source = "jdamata/sonarqube"
    }
  }
}

provider "sonarqube" {
  token = "d4at55a6f7r199bd707h39625685510880gbf7ff"
  host  = "http://127.0.0.1:9000"
}
