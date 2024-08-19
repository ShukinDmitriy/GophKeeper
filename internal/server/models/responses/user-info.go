package responses

type UserInfo struct {
	ID       uint   `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password,omitempty"`
}
