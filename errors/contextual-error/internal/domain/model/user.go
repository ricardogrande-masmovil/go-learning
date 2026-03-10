package model

type User struct {
	id   string
	name string
}

func NewUser(id string, name string) (*User, error) {
	if id == "" {
		return nil, &ValidationError{Cause: CauseInvalidID, Message: "id cannot be empty"}
	}
	if name == "" {
		return nil, &ValidationError{Cause: CauseInvalidName, Message: "name cannot be empty"}
	}

	return &User{
		id:   id,
		name: name,
	}, nil
}

func (u *User) ID() string {
	return u.id
}

func (u *User) Name() string {
	return u.name
}
