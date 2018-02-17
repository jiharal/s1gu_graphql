package question

import (
	"database/sql"

	"github.com/garyburd/redigo/redis"
)

type (
	User struct {
		Id          string
		Username    string
		Password    string
		Email       string
		PhoneNumber string
	}

	AuthUser struct {
		Id          string
		PhoneNumber string
	}
)

func RunCommandRedis(command string, key string, value string) {
	c, err := CachePool.Dial()
	if err != nil {
		Logger.Printf("Error connect to redis: %v", err.Error())
	}
	defer c.Close()
	c.Do(command, key, value)
}

// GetNomorModem get the number modem is available
func GetNomorModem() string {
	for _, x := range PhoneNumberModem {
		data, _ := RedisCheckNumberModemIdle(x)
		if data != "on" {
			return data
			break
		}
		continue
	}
	return ""
}

func GetCommandRedis(command string, key string) string {
	c, err := CachePool.Dial()
	if err != nil {
		Logger.Printf("Error connect to redis: %v", err.Error())
	}
	defer c.Close()
	s, _ := redis.String(c.Do(command, key))

	return s

}

func RedisCheckNumberModemIdle(phonenumber string) (string, error) {
	c, err := CachePool.Dial()
	if err != nil {
		Logger.Printf("Error connect to redis: %v", err.Error())
		return "", err
	}
	defer c.Close()

	s, err := redis.String(c.Do("GET", phonenumber))
	if s != "on" {
		return phonenumber, nil
	}
	return "on", nil

}

// CRUD for redis
func RedisGetPhoneNumberModem(phonenumber string) (string, error) {
	c, err := CachePool.Dial()
	if err != nil {
		Logger.Printf("Error connect to redis: %v", err.Error())
		return "", err
	}
	defer c.Close()

	s, _ := redis.String(c.Do("GET", phonenumber))

	return s, nil
}

// CRUD for cockroachDB

func GetOneUserByPhone(DB *sql.DB, phonenumber string) (*User, error) {
	u := new(User)

	err := DB.QueryRow(`select * from user_authmc where phone_number=$1`, phonenumber).Scan(
		&u.Id,
		&u.Username,
		&u.Password,
		&u.Email,
		&u.PhoneNumber,
	)
	if err == sql.ErrNoRows {
		u.Id = "0"
		u.Username = "Unknown"
		u.Password = "Unknown"
		u.Email = "Unknown"
		u.PhoneNumber = "Unknown"

		return u, nil
	}
	if err != nil {
		Logger.Printf("Error GetOneUserByID : %v", err.Error())
		return nil, err
	}
	return u, nil
}

// GetOneUserById authentication user by id
func GetOneUserById(DB *sql.DB, Username, Password string) (*AuthUser, error) {
	user := new(AuthUser)
	if err := DB.QueryRow(`SELECT id, phone_number from user_authmc WHERE username=$1 and password=$2`, Username, Password).Scan(
		&user.Id,
		&user.PhoneNumber,
	); err == sql.ErrNoRows {
		user.Id = "Unknown"
		user.PhoneNumber = "Unknown"
		return user, err
	}
	return user, nil
}
