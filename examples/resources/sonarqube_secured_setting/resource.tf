resource "sonarqube_secured_setting" "ldap_bind_password" {
  key   = "ldap.bindPassword"
  value = var.ldap_bind_password
}
