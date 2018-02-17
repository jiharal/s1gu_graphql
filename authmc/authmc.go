package authmc

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

/*******************************************************************************************
********************************	TYPE DEFINITIONS	************************************
*******************************************************************************************/

// A AUTHMC is the representation of a AUTHMC GSM modem with several utility features.
type AUTHMC struct {
	port   *SerialPort
	logger *log.Logger
}

/*******************************************************************************************
********************************   GSM: BASIC FUNCTIONS  ***********************************
*******************************************************************************************/

// New creates and initializes a new AUTHMC device.
func New() *AUTHMC {
	return &AUTHMC{
		port:   NewSerial(),
		logger: log.New(os.Stdout, "[AUTHMC] ", log.LstdFlags),
	}
}

// Connect creates a connection with the AUTHMC modem via serial port and test communications.
func (s *AUTHMC) Connect(port string, baud int) error {
	// Open device serial port
	if err := s.port.Open(port, baud, time.Millisecond*100); err != nil {
		return err
	}
	// Ping to Modem
	return s.Ping()
}

func (sim *AUTHMC) Disconnect() error {
	// Close device serial port
	return sim.port.Close()
}

func (sim *AUTHMC) Wait4response(cmd, expected string, timeout time.Duration) (string, error) {
	// Send command
	if err := sim.port.Println(cmd); err != nil {
		return "", err
	}
	// Wait for command response
	regexp := expected + "|" + CMD_ERROR
	response, err := sim.port.WaitForRegexTimeout(regexp, timeout)
	if err != nil {
		return "", err
	}
	// Check if response is an error
	if strings.Contains(response, "ERROR") {
		return response, errors.New("Errors found on command response")
	}
	// Response received succesfully
	return response, nil
}

// Send a SMS
func (s *AUTHMC) SendSMS(number, msg string) error {
	// Set message format
	if err := s.SetSMSMode(TEXT_MODE); err != nil {
		return err
	}
	// Send command
	cmd := fmt.Sprintf(CMD_CMGS, number)
	if err := s.port.Println(cmd); err != nil {
		return err
	}
	// Wait modem to be ready
	time.Sleep(time.Second * 1)
	// Send message
	_, err := s.Wait4response(msg+CMD_CTRL_Z, CMD_OK, time.Second*5)
	if err != nil {
		return err
	}
	// Message sent succesfully
	return nil
}
func (s *AUTHMC) CallBack(number string) error {
	cmd := fmt.Sprintf(CMD_ATD, number+";")
	if err := s.port.Println(cmd); err != nil {
		return err
	}

	ticker := time.NewTicker(time.Millisecond * 500)
	time.Sleep(time.Millisecond * 9000)
	ticker.Stop()
	cmd1 := fmt.Sprintf(CMD_ATH)
	if err := s.port.Println(cmd1); err != nil {
		return err
	}
	// Wait modem to be ready
	time.Sleep(time.Second * 1)
	return nil
}

func (s *AUTHMC) GenerateQuestion(num string) error {
	//slice the variable num from phone number

	user, _ := GetOneUserByPhoneNumber(DbPool, num)

	x := num[8:]
	//convert to string
	x1, _ := strconv.Atoi(x)
	//call randInt function with min and max number to generate from
	x2 := randInt(10, 99)
	// add x1 and x2 and assign it to sum variable
	sum := x1 + x2
	// convert back to string so it can be used as argument in CreateFunction
	x3 := strconv.Itoa(x2)
	//call CreateQuestion that already declare in crud.go
	//that take 4 argument
	CreateQuestion(DbPool, user.Id, "how many sum of "+x+" + "+x3+":", sum)

	return nil
}

//function to generate random int
func randInt(min int, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return min + rand.Intn(max-min)
}

// WaitSMS will return when either a new SMS is recived or the timeout has expired.
// The returned value is the memory ID of the received SMS, use ReadSMS to read SMS content.
func (s *AUTHMC) WaitSMS(timeout time.Duration) (id string, err error) {
	id, err = s.Wait4response("", CMD_CMTI_REGEXP, timeout)
	if err != nil {
		return
	}
	if len(id) >= len(CMD_CMTI_RX) {
		id = id[len(CMD_CMTI_RX):]
	}
	return
}

// ReadSMS retrieves SMS text from inbox memory by ID.
func (s *AUTHMC) ReadSMS(id string) (msg string, err error) {
	// Set message format
	if err := s.SetSMSMode(TEXT_MODE); err != nil {
		return "", err
	}
	// Send command
	cmd := fmt.Sprintf(CMD_CMGR, id)
	if _, err := s.Wait4response(cmd, CMD_CMGR_REGEXP, time.Second*5); err != nil {
		return "", err
	}
	// Reading succesful get message data
	return s.port.ReadLine()
}

// ReadSMS deletes SMS from inbox memory by ID.
func (s *AUTHMC) DeleteSMS(id string) error {
	// Send command
	cmd := fmt.Sprintf(CMD_CMGD, id)
	_, err := s.Wait4response(cmd, CMD_OK, time.Second*1)
	return err
}

// Ping modem
func (s *AUTHMC) Ping() error {
	_, err := s.Wait4response(CMD_AT, CMD_OK, time.Second*1)
	s.Wait4response(CMD_AT_CLIP, CMD_OK, time.Second*1)
	return err
}

// SetSMSMode selects SMS Message Format ("0" = PDU mode, "1" = Text mode)
func (s *AUTHMC) SetSMSMode(mode string) error {
	cmd := fmt.Sprintf(CMD_CMGF_SET, mode)
	_, err := s.Wait4response(cmd, CMD_OK, time.Second*1)
	return err
}

// SetSMSMode reads SMS Message Format (0 = PDU mode, 1 = Text mode)
func (s *AUTHMC) SMSMode() (mode string, err error) {
	mode, err = s.Wait4response(CMD_CMGF, CMD_CMGF_REGEXP, time.Second*1)
	if err != nil {
		return
	}
	if len(mode) >= len(CMD_CMGF_RX) {
		mode = mode[len(CMD_CMGF_RX):]
	}
	return
}

// SetSMSMode selects SMS Message Format (0 = PDU mode, 1 = Text mode)
func (s *AUTHMC) CheckSMSTextMode(mode int) error {
	cmd := fmt.Sprintf(CMD_CMGF, mode)
	_, err := s.Wait4response(cmd, CMD_OK, time.Second*1)
	return err
}
