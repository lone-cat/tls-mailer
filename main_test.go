package tls_mailer_test

import (
	"fmt"
	"github.com/lone-cat/tls-mailer/simpleemail"
	"testing"
)

func TestMy(t *testing.T) {
	email, _ := simpleemail.NewEmptyEmail().WithText(`привет =) aaa`).WithAttachedFile(`./test_attachments/image1.jpg`)
	emailStr := email.String()
	email2, _ := simpleemail.Import(emailStr)
	fmt.Println(email2)
}
