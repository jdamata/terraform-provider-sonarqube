package sonarqube

import "testing"

func TestBadgeTokenSensitivityNextMajor(t *testing.T) {
	for _, tc := range []struct {
		flag string
		want bool
	}{
		{flag: "", want: true},
		{flag: "false", want: true},
		{flag: "true", want: true},
	} {
		t.Run("flag="+tc.flag, func(t *testing.T) {
			t.Setenv("SONARQUBE_BADGE_TOKEN_SENSITIVE", tc.flag)
			provider := Provider()
			for name, resource := range map[string]bool{
				"resource":    provider.ResourcesMap["sonarqube_project"].Schema["badge_token"].Sensitive,
				"data source": provider.DataSourcesMap["sonarqube_project"].Schema["badge_token"].Sensitive,
			} {
				if resource != tc.want {
					t.Errorf("%s badge_token sensitive = %v, want %v", name, resource, tc.want)
				}
			}
		})
	}
}
