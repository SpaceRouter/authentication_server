package forms

type UserLogin struct {
	User     string
	Password string
}

type UserCreate struct {
	User     string
	Password string
	Email    string
	Roles    []string
}

type UserChangeRole struct {
	User string
	Role string
}

type UserChangePassword struct {
	User         string
	PrevPassword string
	NextPassword string
}
