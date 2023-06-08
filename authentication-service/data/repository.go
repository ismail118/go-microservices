package data

type Repository interface {
	GetAll() ([]*User, error)
	GetByEmail(email string) (*User, error)
	GetByID(id int) (*User, error)
	Update(user User) error
	DeleteByID(id int) error
	Insert(user User) (int64, error)
	ResetPassword(password string, user User) error
	PasswordMatches(plainText string, user User) (bool, error)
}
