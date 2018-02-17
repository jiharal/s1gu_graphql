package authmc

import (
	"database/sql"
)

type (
	Question struct {
		Id        string
		IdUser    string
		Question  string
		Answer    string
		Action    string
		Status    bool
		CreatedAt string
		UpdatedAt string
	}

	User struct {
		Id          string
		Username    string
		Password    string
		Email       string
		PhoneNumber string
	}
)

// GET ONE USER
func GetOneUserByPhoneNumber(DB *sql.DB, PhoneNumber string) (*User, error) {
	user := new(User)
	err := DB.QueryRow(`SELECT * from user_authmc WHERE phone_number=$1`, PhoneNumber).Scan(
		&user.Id,
		&user.Username,
		&user.Password,
		&user.Email,
		&user.PhoneNumber,
	)
	if err == sql.ErrNoRows {
		user.Id = "0"
		user.Username = "Unknown"
		user.Password = "Unknown"
		user.PhoneNumber = "Unknonw"
		return user, nil
	}

	if err != nil {
		Logger.Printf("Error GetOneUserByPhoneNumber : %v", err.Error())
		return nil, err
	}
	return user, nil
}

// CREATE QUESTION
func CreateQuestion(DB *sql.DB, iduser, question string, answer int) error {
	_, err := DB.Exec(`INSERT INTO 
		question(iduser, question, answer, action)
		VALUES($1, $2, $3, $4) RETURNING id`, iduser, question, answer, "CREATED")
	//idQuestion, _ := res.LastInsertId()
	if err != nil {
		Logger.Printf("Error insert data into question table: %v", err.Error())
		return err
	}
	return nil
}

// GET QUESTION
func GetQuestionByPhoneNumber(DB *sql.DB, IdUser string) (*Question, error) {
	q := new(Question)
	err := DB.QueryRow(`SELECT * FROM question WHERE iduser=$1 AND action='CREATED'`, IdUser).Scan(
		&q.Id,
		&q.IdUser,
		&q.Question,
		&q.Answer,
		&q.Action,
		&q.Status,
		&q.CreatedAt,
		&q.UpdatedAt,
	)
	if err != nil {
		Logger.Printf("Error GetAnswer : %v", err.Error())
		return nil, err
	}
	return q, nil
}

// GetAnswer get answer from db by id
func GetAnswer(DB *sql.DB, IdQuestion string) (*Question, error) {
	q := new(Question)
	err := DB.QueryRow(`SELECT * FROM question WHERE id=$1 AND action='CREATED'`, IdQuestion).Scan(
		&q.Id,
		&q.IdUser,
		&q.Question,
		&q.Answer,
		&q.Action,
		&q.Status,
		&q.CreatedAt,
		&q.UpdatedAt,
	)
	if err != nil {
		Logger.Printf("Error GetAnswer : %v", err.Error())
		return nil, err
	}
	return q, nil
}
