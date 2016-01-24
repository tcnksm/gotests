package method

type User struct {
	Name, Password string
}

func (u *User) update() error {
}

func (u *User) Validate() error {
	return nil
}
