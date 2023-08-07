package user

import (
	"database/sql"
	"github.com/evermos/boilerplate-go/infras"
	"github.com/evermos/boilerplate-go/shared/failure"
	"github.com/evermos/boilerplate-go/shared/logger"
	"github.com/gofrs/uuid"
	"github.com/rs/zerolog/log"
	"regexp"
)

var (
	usersQueries = struct {
		selectUsers string
		insertUsers string
	}{
		insertUsers: `
			INSERT INTO users(
				user_id,
			    username,
			    email,
			    password,
			    role
			) VALUES (
				:user_id,
			    :username,
			    :email,
			    :password,
			    :role)`,
		selectUsers: `
			SELECT
			    u.user_id,
			    u.username,
			    u.email,
			    u.password,
			    u.role
			FROM users u`,
	}
)

type UserRepository interface {
	Create(user Users) (err error)
	ExistsByID(id uuid.UUID) (exists bool, err error)
	ResolveByEmail(email string) (user Users, err error)
}

type UserRepositoryMySQL struct {
	DB *infras.MySQLConn
}

func ProvideUserRepositoryMysql(db *infras.MySQLConn) *UserRepositoryMySQL {
	return &UserRepositoryMySQL{DB: db}
}

func (u *UserRepositoryMySQL) Create(user Users) (err error) {
	exists, err := u.ExistsByID(user.ID)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	if exists {
		err = failure.NotFound("users")
		logger.ErrorWithStack(err)
		return
	}
	isAvailble, err := u.checkEmail(user.Email)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	if !isAvailble {
		log.Info().Msg("Email has been used")
		return
	}

	if !isValidEmail(user.Email) {
		log.Info().Msg("Invalid Email")
		return
	}

	err = u.insertUser(user)
	if err != nil {
		logger.ErrorWithStack(err)
	}

	return
}

func (u *UserRepositoryMySQL) ExistsByID(id uuid.UUID) (exists bool, err error) {
	err = u.DB.Read.Get(
		&exists,
		"SELECT COUNT(user_id) FROM users u WHERE u.user_id = ?",
		id.String())
	if err != nil {
		logger.ErrorWithStack(err)
	}

	return
}

func (u *UserRepositoryMySQL) ResolveByEmail(email string) (user Users, err error) {
	err = u.DB.Read.Get(
		&user,
		usersQueries.selectUsers+" WHERE u.email = ?", email)
	if err != nil && err == sql.ErrNoRows {
		err = failure.NotFound("users")
		logger.ErrorWithStack(err)
		return
	}
	return
}

func (u *UserRepositoryMySQL) checkEmail(email string) (bool, error) {
	var count int
	err := u.DB.Read.Get(&count, "SELECT COUNT(*) FROM users WHERE email = ?", email)
	if err != nil {
		return false, err
	}
	if count > 0 {
		return false, err
	}
	return true, err
}

func isValidEmail(email string) bool {
	regex, err := regexp.Compile("^[\\w-\\.]+@([\\w-]+\\.)+[\\w-]{2,4}$")
	if err != nil {
		return false
	}
	return regex.MatchString(email)
}

func (u *UserRepositoryMySQL) insertUser(user Users) error {
	stmt, err := u.DB.Write.PrepareNamed(usersQueries.insertUsers)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(user)
	return err
}
