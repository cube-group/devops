package models

func GitlabOAuthURL(ref string) (res string) {
	gitlab, err := Gitlab()
	if err != nil {
		return
	}
	return gitlab.GetAuthURL(ref)
}
