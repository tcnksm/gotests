package method

type User struct {
	Name, Password string
}

func NewUser(s string) *User {
	return nil
}

func (u *User) initialize() error {
}

func (u User) Add() error {
	return nil
}

func (u *User) Delete(name string) error {
	return nil
}
