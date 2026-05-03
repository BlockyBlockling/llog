package llog

import (
	"os"
	"strconv"
	"testing"

	"github.com/joho/godotenv"
)

func TestMail(t *testing.T) {
	//Running Mail Test
	t.Log("Running Mail Test:")

	// get .env file
	err := godotenv.Load()
	if err.Error() == "open .env: no such file or directory" {
		// continue
	} else if err != nil {
		t.Error(err)
	}

	TEST_MAIL_ADDRESS := os.Getenv("TEST_MAIL_ADDRESS")
	TEST_MAIL_TARGET := os.Getenv("TEST_MAIL_TARGET")
	TEST_MAIL_URL := os.Getenv("TEST_MAIL_URL")
	TEST_MAIL_PORT_STRING := os.Getenv("TEST_MAIL_PORT")
	TEST_MAIL_PORT, err := strconv.Atoi(TEST_MAIL_PORT_STRING)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	TEST_MAIL_USERNAME := os.Getenv("TEST_MAIL_USERNAME")
	TEST_MAIL_PASSWORD := os.Getenv("TEST_MAIL_PASSWORD")

	InitMail(TEST_MAIL_ADDRESS, []string{TEST_MAIL_TARGET}, TEST_MAIL_URL, TEST_MAIL_PORT, TEST_MAIL_USERNAME, TEST_MAIL_PASSWORD, "Automatic llog Testing")
	err = NotifyMail("Test Notify")
	if err != nil {
		t.Error("Mail Send Error: " + err.Error())
		t.Fail()
	}
}
