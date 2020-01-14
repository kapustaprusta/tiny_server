package defs

type Endpoint struct {
	URL         string `json:"url"`
	Method      string `json:"method"`
	Description string `json:"description"`
}

type API struct {
	Endpoints []Endpoint `json:"endpoints"`
}

type User struct {
	ID        int    `json:"id"`
	Nickname  string `json:"username"`
	Email     string `json:"email"`
	FirstName string `json:"first"`
	LastName  string `json:"last"`
}

type Users struct {
	Users []User `json:"users"`
}

type CreationResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ErrorMessage struct {
	Code     int
	HTTPCode int
	Comment  string
}
