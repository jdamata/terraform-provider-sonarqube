# Import global new code period
terraform import sonarqube_new_code_periods.example newCodePeriod

# Import project-specific new code period
terraform import sonarqube_new_code_periods.example newCodePeriod/my-project-key

# Import branch-specific new code period
terraform import sonarqube_new_code_periods.example newCodePeriod/my-branch-name/my-project-key
