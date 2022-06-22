package tls_mailer_test

import (
	"fmt"
	"github.com/lone-cat/tls-mailer/simpleemail"
	"testing"
)

func TestMy(t *testing.T) {
	email := simpleemail.NewEmptyEmail()
	fmt.Println(email)
}
