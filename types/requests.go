package types

type CreateEmailRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type SucessResponse struct {
	Success bool `json:"success"`
}

type ListEmailsResponse struct {
	Success bool    `json:"success"`
	Emails  []Email `json:"emails"`
}
