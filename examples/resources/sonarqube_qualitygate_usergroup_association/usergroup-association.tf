resource "sonarqube_qualitygate" "main" {
  name = "my_qualitygate"

  condition {
    metric    = "new_coverage"
    op        = "LT"
    threshold = "30"
  }
}

resource "sonarqube_group" "qa_team" {
  name        = "QA-Team"
  description = "Quality Assurence Team"
}

resource "sonarqube_qualitygate_usergroup_association" "main" {
  gatename   = sonarqube_qualitygate.main.id
  group_name = sonarqube_group.qa_team.name
}
