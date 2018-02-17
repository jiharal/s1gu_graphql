package question

import (
	"context"
	"fmt"
	"log"
	"time"
)

// Modem number
var (
	PhoneNumberModem = []string{"0895421313771", "082325600997", "082325600998"}
	ExpireTime       = "60"
)

type Resolver struct{}

type (
	authmessage struct {
		Message          string
		Id               string
		Phonenumber      string
		Phonenumbermodem string
	}
	randomQuestion struct {
		Token    string
		Question string
	}

	questionResolver struct {
		q *randomQuestion
	}

	// auth user struct

	// Authuser resolver
	authResolver struct {
		a *authmessage
	}
)

var (
	questionData = make(map[string]*randomQuestion)
)

func (r *Resolver) RandomQuestion(args struct{ Token string }) *questionResolver {
	startTime := time.Now()
	x, y := 2, 6

	var a = &randomQuestion{
		Token:    "12345",
		Question: fmt.Sprintf("Berapa hasil penjumlahan dari %d + %d :", x, y),
	}
	questionData[a.Token] = a

	if a := questionData[args.Token]; a != nil {
		endTime := time.Now()
		log.Println("duration Graphql:", endTime.Sub(startTime).Seconds())
		return &questionResolver{a}
	}
	return nil
}

func (r *questionResolver) Token() string {
	return r.q.Token
}

func (r *questionResolver) Question() string {
	return r.q.Question
}

// Authuser flow 1 - 21
func (r *Resolver) AuthUser(ctx context.Context, args struct {
	Iduser      string
	Username    string
	Password    string
	PhoneNumber string
}) *authResolver {
	a, _ := GetOneUserById(DbPool, args.Username, args.Password)
	if a.Id == args.Iduser {
		if a.Id == args.Iduser && a.PhoneNumber == args.PhoneNumber {
			PhoneNumberOnTemp := GetNomorModem()
			fmt.Println("Bisa di pake:", PhoneNumberOnTemp)

			// Set phone number modem = on
			// this is use to check if this number is available
			RunCommandRedis("SET", PhoneNumberOnTemp, "on")

			// Set phone number = phone number modem
			// this is use to get modem number by number user for delete modem number on redis
			RunCommandRedis("SET", a.PhoneNumber, PhoneNumberOnTemp)

			// Set modem phone number = phone number modem
			// This is use to get number user in redis for to callback user.
			RunCommandRedis("SET", "authmc-"+PhoneNumberOnTemp, a.PhoneNumber)

			// Set expire time.
			RunCommandRedis("EXPIRE", PhoneNumberOnTemp, ExpireTime)
			RunCommandRedis("EXPIRE", a.PhoneNumber, ExpireTime)
			RunCommandRedis("EXPIRE", "authmc-"+PhoneNumberOnTemp, ExpireTime)

			m := &authmessage{
				Message:          "Please call to number : " + PhoneNumberOnTemp,
				Id:               a.Id,
				Phonenumber:      a.PhoneNumber,
				Phonenumbermodem: PhoneNumberOnTemp,
			}
			return &authResolver{m}
		} else {
			m := &authmessage{
				Message: "Phone Number Not Found!",
			}
			return &authResolver{m}
		}

	} else {
		m := &authmessage{
			Message: "ID not found!",
		}
		return &authResolver{m}
	}
	return nil
}

func (r *authResolver) Message() string {
	return r.a.Message
}

func (r *authResolver) Id() string {
	return r.a.Id
}

func (r *authResolver) Phonenumber() string {
	return r.a.Phonenumber
}

func (r *authResolver) Phonenumbermodem() string {
	return r.a.Phonenumbermodem
}
