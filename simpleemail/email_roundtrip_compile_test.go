package simpleemail

import (
	"fmt"
	"testing"
)

func TestRoundtripCompile(t *testing.T) {
	for _, em := range emailsForTest {
		testComiledCompare(em.(*email), t)
	}
}

func testComiledCompare(email *email, t *testing.T) {
	emailString := email.String()

	importedEmail, err := Import(emailString)
	if err != nil {
		fmt.Println(`import failed:`)
		fmt.Println(emailString)

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
