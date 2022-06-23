package simpleemail_test

import (
	"fmt"
	"github.com/lone-cat/tls-mailer/simpleemail"
	"testing"
)

func TestHeadersOnly(t *testing.T) {
	email := simpleemail.NewEmptyEmail()

	emailsForTest := make([]simpleemail.Email, 0)

	emailsForTest = append(emailsForTest, email)
	email = email.
		WithFrom(from)
	emailsForTest = append(emailsForTest, email)
	email = email.
		WithTo(to)
	emailsForTest = append(emailsForTest, email)
	email = email.
		WithCc(cc)
	emailsForTest = append(emailsForTest, email)
	email = email.
		WithBcc(bcc)
	emailsForTest = append(emailsForTest, email)
	email = email.
		WithSubject(subject)
	emailsForTest = append(emailsForTest, email)
	email2 := email.
		WithText(text)
	emailsForTest = append(emailsForTest, email2)
	email2 = email.
		WithHtml(html)
	emailsForTest = append(emailsForTest, email2)
	email2 = email.
		WithEmbeddedString(embedded)
	emailsForTest = append(emailsForTest, email2)
	email2 = email.
		WithAttachedString(attached)
	emailsForTest = append(emailsForTest, email2)
	email2 = email.
		WithText(text)
	email3 := email2.
		WithHtml(html)
	emailsForTest = append(emailsForTest, email3)
	email3 = email2.
		WithEmbeddedString(embedded)
	emailsForTest = append(emailsForTest, email3)
	email3 = email2.
		WithAttachedString(attached)
	emailsForTest = append(emailsForTest, email3)
	email2 = email.
		WithHtml(html)
	email3 = email2.
		WithEmbeddedString(embedded)
	emailsForTest = append(emailsForTest, email3)
	email3 = email2.
		WithAttachedString(attached)
	emailsForTest = append(emailsForTest, email3)
	email2 = email.
		WithEmbeddedString(embedded)
	email3 = email2.
		WithAttachedString(attached)
	emailsForTest = append(emailsForTest, email3)
	email3 = email.
		WithText(text).
		WithHtml(html).
		WithEmbeddedString(embedded).
		WithAttachedString(attached)
	emailsForTest = append(emailsForTest, email3)

	for _, em := range emailsForTest {
		testEmail(em, t)
	}
}

func testEmail(email simpleemail.Email, t *testing.T) {
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
