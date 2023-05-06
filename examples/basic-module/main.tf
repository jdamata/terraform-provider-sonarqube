terraform {
  required_providers {
    sonarqube = {
      source = "jdamata/sonarqube"
    }
  }
}

provider "sonarqube" {
    user   = "admin"
    pass   = "admin" 
    host   = "http://127.0.0.1:9000"
}

resource "sonarqube_qualitygate" "tf-postfix" {
	name = "tf-postfix"

  condition {
    metric    = "new_coverage"
    op        = "LT"
    threshold = "30"
  }
}

module "tf-postfix-repo" {
    source   = "./modules/example"
    
    name                 = "tf-postfix"
    quality_gates = {
        gate1 = {
            id        = sonarqube_qualitygate.tf-postfix.id
            metric    = "vulnerabilities"
            threshold = "11"
            operator  = "GT"
        }
    }
}