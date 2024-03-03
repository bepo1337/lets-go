package models

import (
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)

const (
	UserHashCost = 12
)

type User struct {
	Id       int
	Name     string
	Email    string
	HashedPw []byte
	Created  time.Time
}

type UserModelInterface interface {
	Insert(name, email, password string) error
	Exists(id int) (bool, error)
	Authenticate(email, password string) (int, error)
	Get(id int) (*User, error)
	CorrectPassword(id int, password string) (bool, error)
	UpdatePassword(id int, newPassword string) (bool, error)
}

type UserModel struct {
	DB *sql.DB
}

func (u *UserModel) Insert(name, email, password string) error {
	hashedPw, err := bcrypt.GenerateFromPassword([]byte(password), UserHashCost)
	if err != nil {
		return err
	}
	dbStatement := `INSERT INTO users (name, email, hashed_pw, created)
						values(?, ?, ?, UTC_TIMESTAMP())`
	_, err = u.DB.Exec(dbStatement, name, email, hashedPw)
	if err != nil {
		var mySqlError *mysql.MySQLError
		if errors.As(err, &mySqlError) {
			if mySqlError.Number == 1062 && strings.Contains(mySqlError.Message, "users_uc_email") {
				return ErrDuplicateEmail
			}
		}
		return err
	}
	return nil
}

func (u *UserModel) Exists(id int) (bool, error) {
	dbStatement := `SELECT id FROM users WHERE id=?`
	row := u.DB.QueryRow(dbStatement, id)
	var test int
	err := row.Scan(&test)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (u *UserModel) Authenticate(email, password string) (int, error) {
	dbStatement := "SELECT id, hashed_pw  FROM users where email=?"
	var user = &User{}
	row := u.DB.QueryRow(dbStatement, email)
	err := row.Scan(&user.Id, &user.HashedPw)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return -1, ErrInvalidCredentials
		}
		return -1, err
	}
	err = bcrypt.CompareHashAndPassword(user.HashedPw, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return -1, ErrInvalidCredentials
		}
		return -1, err
	}
	return user.Id, nil
}

func (u *UserModel) Get(id int) (*User, error) {
	dbStatement := "SELECT id, email, created, name  FROM users where id=?"
	var user = &User{}
	row := u.DB.QueryRow(dbStatement, id)
	err := row.Scan(&user.Id, &user.Email, &user.Created, &user.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	return user, nil
}

func (u *UserModel) CorrectPassword(id int, password string) (bool, error) {
	dbStatement := "SELECT hashed_pw  FROM users where id=?"
	//dbStatement := "SELECT EXISTS(SELECT 1 FROM users where id=? and hashed_pw=?)"
	//hashedPw, err := bcrypt.GenerateFromPassword([]byte(password), UserHashCost)
	//if err != nil {
	//	return false, err
	//}
	var hashedPassword string
	err := u.DB.QueryRow(dbStatement, id).Scan(&hashedPassword)
	if err != nil {
		return false, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return false, ErrHashesDontMatch
	}
	return true, nil
}

func (u *UserModel) UpdatePassword(id int, newPassword string) (bool, error) {
	dbStatement := "UPDATE users set hashed_pw=? where id=?"
	newHashedPw, err := bcrypt.GenerateFromPassword([]byte(newPassword), UserHashCost)
	if err != nil {
		return false, err
	}
	result, err := u.DB.Exec(dbStatement, newHashedPw, id)
	if err != nil {
		return false, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	if rowsAffected == 1 {
		return true, nil
	} else {
		return false, nil
	}
}
