package data

import (
	"context"
	"database/sql"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"time"
)

const dbTimeout = time.Second * 3

var db *sql.DB

type Models struct {
	User User
}

type User struct {
	ID        int
	Email     string
	FirstName string    `json:"first_name,omitempty"`
	LastName  string    `json:"last_name,omitempty"`
	Password  string    `json:"-"`
	Active    int       `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func New(dbPool *sql.DB) Models {
	db = dbPool

	return Models{
		User: User{},
	}
}

func (u *User) GetAll() ([]*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, email, first_name, last_name, password, user_active, created_at, updated_at 
	from users order by last_name`

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User

	for rows.Next() {
		var x User
		err = rows.Scan(
			&x.ID,
			&x.FirstName,
			&x.LastName,
			&x.Password,
			&x.Active,
			&x.CreatedAt,
			&x.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, &x)
	}

	err = rows.Err()
	if err != nil {
		return users, err
	}

	return users, nil
}

func (u *User) GetByEmail(email string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, email, first_name, last_name, password, user_active, created_at, updated_at 
	from users where email = $1`

	row := db.QueryRowContext(ctx, query, email)

	var user User

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.Active,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	err = row.Err()
	if err != nil {
		return &user, err
	}

	return &user, nil
}

func (u *User) GetByID(id int) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, email, first_name, last_name, password, user_active, created_at, updated_at 
	from users where id = $1`

	row := db.QueryRowContext(ctx, query, id)

	var user User

	err := row.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.Active,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	err = row.Err()
	if err != nil {
		return &user, err
	}

	return &user, nil
}

func (u *User) Update() error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `update users set 
                 email = $1,
                 first_name = $2,
                 last_name = $3,
                 user_active = $4,
                 updated_at = $5 
             where id = $6`

	_, err := db.ExecContext(ctx, query,
		u.Email,
		u.FirstName,
		u.LastName,
		u.Active,
		time.Now(),
		u.ID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (u *User) Delete() error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `delete from users where id = $1`

	_, err := db.ExecContext(ctx, query, u.ID)
	if err != nil {
		return err
	}

	return nil
}

func (u *User) DeleteByID(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `delete from users where id = $1`

	_, err := db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}

func (u *User) Insert(user User) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)
	if err != nil {
		return 0, err
	}

	var newID int64
	query := `insert into users (email, first_name, last_name, password, user_active, created_at, updated_at) 
	values ($1, $2, $3, $4, $5, $6, $7)`

	res, err := db.ExecContext(ctx, query,
		u.Email,
		u.FirstName,
		u.LastName,
		hashedPassword,
		u.Active,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return 0, err
	}

	newID, err = res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return newID, nil
}

func (u *User) ResetPassword(password string) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return err
	}

	query := `update users set password = $1 where id = $2`

	_, err = db.ExecContext(ctx, query, hashedPassword, u.ID)
	if err != nil {
		return err
	}

	return nil
}

func (u *User) PasswordMatches(plainText string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plainText))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			// invalid password
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}
