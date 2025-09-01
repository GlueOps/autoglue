package orgs

type OrgInput struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type InviteInput struct {
	Email string `json:"email"`
	Role  string `json:"role"`
}
