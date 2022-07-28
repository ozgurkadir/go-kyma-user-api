package user

type User struct {
	UserName  string `json:"username`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Address   string `json:"address"`
	Mobile    int    `json:"mobile"`
}
