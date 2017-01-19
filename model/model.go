package model

type User struct {
	ID       string `json:"id"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type Child struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Parent struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Children []string `json:"children"`
}

type Project struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}