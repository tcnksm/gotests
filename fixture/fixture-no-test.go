package fixture

func DoSomething() error {
	return nil
}

type User struct {
	Name, Password string
}

func (u *User) Validate() error {
	return nil
}
