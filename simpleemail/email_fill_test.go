package simpleemail_test

import (
	"fmt"
	"github.com/lone-cat/tls-mailer/simpleemail"
	"os"
	"testing"
)

func TestEmpty(t *testing.T) {
	email := simpleemail.NewEmptyEmail()
	if email.GetText() != `` {
		fmt.Printf("`Text` is `%s` in empty email\r\n", email.GetText())
		t.Fail()
	}
	if email.GetHtml() != `` {
		fmt.Printf("`Html` is `%s` in empty email\r\n", email.GetHtml())
		t.Fail()
	}
	if email.GetSubject() != `` {
		fmt.Printf("`Subject` is `%s` in empty email\r\n", email.GetSubject())
		t.Fail()
	}
	if len(email.GetFrom()) > 0 {
		fmt.Printf("`From` contains %#v instead of empty list\r\n", email.GetFrom())
		t.Fail()
	}
	if len(email.GetTo()) > 0 {
		fmt.Printf("`To` contains %#v instead of empty list\r\n", email.GetTo())
		t.Fail()
	}
	if len(email.GetCc()) > 0 {
		fmt.Printf("`Cc` contains %#v instead of empty list\r\n", email.GetCc())
		t.Fail()
	}
	if len(email.GetBcc()) > 0 {
		fmt.Printf("`Bcc` contains %#v instead of empty list\r\n", email.GetBcc())
		t.Fail()
	}
	if len(email.GetEmbedded()) > 0 {
		fmt.Println("`Embedded` list is not empty in empty email")
		t.Fail()
	}
	if len(email.GetAttachments()) > 0 {
		fmt.Println("`Attachments` list is not empty in empty email")
		t.Fail()
	}
}

func TestFill(t *testing.T) {
	email := simpleemail.NewEmptyEmail().
		WithText(text).
		WithHtml(html).
		WithSubject(subject).
		WithFrom(from).
		WithTo(to).
		WithCc(cc).
		WithBcc(bcc)
	email, err := email.
		WithEmbeddedFile(`cid`, embedded)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}

	email, err = email.
		WithAttachedFile(attached)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}

	if email.GetText() != text {
		fmt.Printf("`Text` filled incorrectly. expected %#v, got %#v\r\n", text, email.GetText())
		t.Fail()
	}
	if email.GetHtml() != html {
		fmt.Printf("`Html` filled incorrectly. expected %#v, got %#v\r\n", html, email.GetHtml())
		t.Fail()
	}
	if email.GetSubject() != subject {
		fmt.Printf("`Subject` filled incorrectly. expected %#v, got %#v\r\n", subject, email.GetSubject())
		t.Fail()
	}
	if !addressSlicesEqual(email.GetFrom(), from) {
		fmt.Printf("`From` filled incorrectly. expected %#v, got %#v\r\n", from, email.GetFrom())
		t.Fail()
	}
	if !addressSlicesEqual(email.GetTo(), to) {
		fmt.Printf("`To` filled incorrectly. expected %#v, got %#v\r\n", to, email.GetTo())
		t.Fail()
	}
	if !addressSlicesEqual(email.GetCc(), cc) {
		fmt.Printf("`Cc` filled incorrectly. expected %#v, got %#v\r\n", cc, email.GetCc())
		t.Fail()
	}
	if !addressSlicesEqual(email.GetBcc(), bcc) {
		fmt.Printf("`Bcc` filled incorrectly. expected %#v, got %#v\r\n", bcc, email.GetBcc())
		t.Fail()
	}
	dataBytes, err := os.ReadFile(embedded)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	if email.GetEmbedded()[0].GetBody() != string(dataBytes) {
		fmt.Printf("`Embedded body` filled incorrectly. expected %#v, got %#v\r\n", string(dataBytes), email.GetEmbedded()[0].GetBody())
		t.Fail()
	}
	dataBytes, err = os.ReadFile(attached)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	if email.GetAttachments()[0].GetBody() != string(dataBytes) {
		fmt.Printf("`Attached body` filled incorrectly. expected %#v, got %#v\r\n", string(dataBytes), email.GetAttachments()[0].GetBody())
		t.Fail()
	}
}
