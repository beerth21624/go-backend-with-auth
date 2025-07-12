package domain

import (
	"time"
)

type User struct {
	id        UserID
	username  NonEmptyString
	email     NonEmptyString
	firstName NonEmptyString
	lastName  NonEmptyString
	password  HashedPassword
	role      UserRole
	status    Status
	createdAt CreatedAt
	updatedAt UpdatedAt
}

func NewUser(
	username, email, firstName, lastName, plainPassword string,
) (*User, error) {
	usernameVO, err := NewNonEmptyString(username)
	if err != nil {
		return nil, err
	}

	emailVO, err := NewNonEmptyString(email)
	if err != nil {
		return nil, err
	}

	firstNameVO, err := NewNonEmptyString(firstName)
	if err != nil {
		return nil, err
	}

	lastNameVO, err := NewNonEmptyString(lastName)
	if err != nil {
		return nil, err
	}

	hashedPassword, err := NewHashedPassword(plainPassword)
	if err != nil {
		return nil, err
	}

	role, err := NewUserRole("user")
	if err != nil {
		return nil, err
	}

	status, err := NewStatus("active")
	if err != nil {
		return nil, err
	}

	now := time.Now()
	createdAt, err := NewCreatedAt(now)
	if err != nil {
		return nil, err
	}

	updatedAt, err := NewUpdatedAt(now)
	if err != nil {
		return nil, err
	}

	return &User{
		id:        NewUserID(),
		username:  usernameVO,
		email:     emailVO,
		firstName: firstNameVO,
		lastName:  lastNameVO,
		password:  hashedPassword,
		role:      role,
		status:    status,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}, nil
}

func ReconstructUser(
	id, username, email, firstName, lastName, hashedPassword, role, status string,
	createdAt, updatedAt time.Time,
) (*User, error) {
	idVO, err := NewUserIDFromString(id)
	if err != nil {
		return nil, err
	}

	usernameVO, err := NewNonEmptyString(username)
	if err != nil {
		return nil, err
	}

	emailVO, err := NewNonEmptyString(email)
	if err != nil {
		return nil, err
	}

	firstNameVO, err := NewNonEmptyString(firstName)
	if err != nil {
		return nil, err
	}

	lastNameVO, err := NewNonEmptyString(lastName)
	if err != nil {
		return nil, err
	}

	passwordVO, err := NewHashedPasswordFromHash(hashedPassword)
	if err != nil {
		return nil, err
	}

	roleVO, err := NewUserRole(role)
	if err != nil {
		return nil, err
	}

	statusVO, err := NewStatus(status)
	if err != nil {
		return nil, err
	}

	createdAtVO, err := NewCreatedAt(createdAt)
	if err != nil {
		return nil, err
	}

	updatedAtVO, err := NewUpdatedAt(updatedAt)
	if err != nil {
		return nil, err
	}

	return &User{
		id:        idVO,
		username:  usernameVO,
		email:     emailVO,
		firstName: firstNameVO,
		lastName:  lastNameVO,
		password:  passwordVO,
		role:      roleVO,
		status:    statusVO,
		createdAt: createdAtVO,
		updatedAt: updatedAtVO,
	}, nil
}

func (u *User) ID() UserID {
	return u.id
}

func (u *User) Username() NonEmptyString {
	return u.username
}

func (u *User) Email() NonEmptyString {
	return u.email
}

func (u *User) FirstName() NonEmptyString {
	return u.firstName
}

func (u *User) LastName() NonEmptyString {
	return u.lastName
}

func (u *User) Password() HashedPassword {
	return u.password
}

func (u *User) FullName() string {
	return u.firstName.String() + " " + u.lastName.String()
}

func (u *User) Role() UserRole {
	return u.role
}

func (u *User) Status() Status {
	return u.status
}

func (u *User) CreatedAt() CreatedAt {
	return u.createdAt
}

func (u *User) UpdatedAt() UpdatedAt {
	return u.updatedAt
}

func (u *User) IsActive() bool {
	return u.status.IsActive()
}

func (u *User) Activate() error {
	status, err := NewStatus("active")
	if err != nil {
		return err
	}
	u.status = status
	u.updatedAt = NewUpdatedAtNow()
	return nil
}

func (u *User) Deactivate() error {
	status, err := NewStatus("inactive")
	if err != nil {
		return err
	}
	u.status = status
	u.updatedAt = NewUpdatedAtNow()
	return nil
}

func (u *User) UpdateProfile(firstName, lastName string) error {
	firstNameVO, err := NewNonEmptyString(firstName)
	if err != nil {
		return err
	}

	lastNameVO, err := NewNonEmptyString(lastName)
	if err != nil {
		return err
	}

	u.firstName = firstNameVO
	u.lastName = lastNameVO
	u.updatedAt = NewUpdatedAtNow()
	return nil
}

func (u *User) VerifyPassword(plainPassword string) bool {
	return u.password.VerifyPassword(plainPassword)
}

func (u *User) ChangePassword(newPlainPassword string) error {
	newHashedPassword, err := NewHashedPassword(newPlainPassword)
	if err != nil {
		return err
	}

	u.password = newHashedPassword
	u.updatedAt = NewUpdatedAtNow()
	return nil
}

func (u *User) CanLogin() bool {
	return u.IsActive()
}
