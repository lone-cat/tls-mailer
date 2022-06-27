package simpleemail_test

import (
	"fmt"
	"github.com/lone-cat/tls-mailer/simpleemail"
	"testing"
)

func TestRoundtripCompile(t *testing.T) {
	for _, em := range emailsForTest {
		testComiledCompare(em, t)
	}
}

func testComiledCompare(email *simpleemail.Email, t *testing.T) {
	emailString := email.String()

	importedEmail, err := simpleemail.Import(emailString)
	if err != nil {
		fmt.Println(`import failed:`)
		fmt.Println(err)
		t.FailNow()
	}

	importedEmailString := importedEmail.String()
	if emailString != importedEmailString {
		fmt.Println(`twice converted email does not match:`)
		fmt.Println(`--- original ---`)
		fmt.Println(emailString)
		fmt.Println(`--- reconstructed ---`)
		fmt.Println(importedEmailString)
		fmt.Println(`--- end ---`)
		t.Fail()
	}
}
