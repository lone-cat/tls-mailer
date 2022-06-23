package simpleemail_test

import (
	"fmt"
	"github.com/lone-cat/tls-mailer/simpleemail"
	"net/mail"
	"testing"
)

func addressSlicesEqual(addressSlice1 []mail.Address, addressSlice2 []mail.Address) bool {
	addressMap1 := convertAddressSliceToMapByEmail(addressSlice1)
	addressMap2 := convertAddressSliceToMapByEmail(addressSlice2)
	if len(addressMap1) != len(addressMap2) {
		return false
	}

	var ok bool
	var addr2 mail.Address
	for email, addr := range addressMap1 {
		addr2, ok = addressMap2[email]
		if !ok {
			return false
		}
		if addr.Name != addr2.Name {
			return false
		}
		if addr.Address != addr2.Address {
			return false
		}
	}

	return true
}

func convertAddressSliceToMapByEmail(addressSlice []mail.Address) map[string]mail.Address {
	result := make(map[string]mail.Address)
	for _, addr := range addressSlice {
		result[addr.Address] = mail.Address{Name: addr.Name, Address: addr.Name}
	}
	return result
}

var (
	emptyEmail = simpleemail.NewEmptyEmail().
		WithFrom(from).
		WithTo(to).
		WithCc(cc).
		WithBcc(bcc).
		WithSubject(subject)
)

func TestRoundTripHeaders(t *testing.T) {
	email := emptyEmail

	emailString := email.String()

	importedEmail, err := simpleemail.Import(emailString)
	if err != nil {
		fmt.Println(`import failed:`)
		fmt.Println(err)
		t.FailNow()
	}

	importedFrom := importedEmail.GetFrom()
	if !addressSlicesEqual(from, importedFrom) {
		fmt.Printf("`From` corrupted. expected %#v, got %#v\r\n", from, importedFrom)
		t.Fail()
	}

	importedTo := importedEmail.GetTo()
	if !addressSlicesEqual(to, importedTo) {
		fmt.Printf("`To` corrupted. expected %#v, got %#v\r\n", to, importedTo)
		t.Fail()
	}

	importedCc := importedEmail.GetCc()
	if !addressSlicesEqual(cc, importedCc) {
		fmt.Printf("`Cc` corrupted. expected %#v, got %#v\r\n", cc, importedCc)
		t.Fail()
	}

	importedBcc := importedEmail.GetBcc()
	if !addressSlicesEqual(bcc, importedBcc) {
		fmt.Printf("`Bcc` corrupted. expected %#v, got %#v\r\n", bcc, importedBcc)
		t.Fail()
	}

	importedSubject := importedEmail.GetSubject()
	if importedSubject != subject {
		fmt.Printf("`Subject` corrupted. expected \"%s\", got \"%s\"\r\n", subject, importedSubject)
		t.Fail()
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

func TestRoundTripText(t *testing.T) {
	email := emptyEmail.
		WithText(text)

	emailString := email.String()

	importedEmail, err := simpleemail.Import(emailString)
	if err != nil {
		fmt.Println(`import failed:`)
		fmt.Println(err)
		t.FailNow()
	}

	importedText := importedEmail.GetText()
	if importedText != text {
		fmt.Printf("`text` corrupted. expected \"%s\", got \"%s\"\r\n", text, importedText)
		t.Fail()
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

func TestRoundTripHtml(t *testing.T) {
	email := emptyEmail.
		WithHtml(html)

	emailString := email.String()

	importedEmail, err := simpleemail.Import(emailString)
	if err != nil {
		fmt.Println(`import failed:`)
		fmt.Println(err)
		t.FailNow()
	}

	importedHtml := importedEmail.GetHtml()
	if importedHtml != html {
		fmt.Printf("`html` corrupted. expected \"%s\", got \"%s\"\r\n", html, importedHtml)
		t.Fail()
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

func TestRoundTripTextAndHtml(t *testing.T) {
	email := emptyEmail.
		WithText(text).
		WithHtml(html)

	emailString := email.String()

	importedEmail, err := simpleemail.Import(emailString)
	if err != nil {
		fmt.Println(`import failed:`)
		fmt.Println(err)
		t.FailNow()
	}

	importedText := importedEmail.GetText()
	if importedText != text {
		fmt.Printf("`text` corrupted. expected \"%s\", got \"%s\"\r\n", text, importedText)
		t.Fail()
	}

	importedHtml := importedEmail.GetHtml()
	if importedHtml != html {
		fmt.Printf("`html` corrupted. expected \"%s\", got \"%s\"\r\n", html, importedHtml)
		t.Fail()
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

func TestRoundTripTextAndHtmlAndEmbedded(t *testing.T) {
	email := emptyEmail.
		WithText(text).
		WithHtml(html).
		WithEmbeddedString(embedded)

	emailString := email.String()

	importedEmail, err := simpleemail.Import(emailString)
	if err != nil {
		fmt.Println(`import failed:`)
		fmt.Println(err)
		t.FailNow()
	}

	importedText := importedEmail.GetText()
	if importedText != text {
		fmt.Printf("`text` corrupted. expected \"%s\", got \"%s\"\r\n", text, importedText)
		t.Fail()
	}

	importedHtml := importedEmail.GetHtml()
	if importedHtml != html {
		fmt.Printf("`html` corrupted. expected \"%s\", got \"%s\"\r\n", html, importedHtml)
		t.Fail()
	}

	importedEmbeddedParts := importedEmail.GetEmbedded()
	if len(importedEmbeddedParts) < 1 {
		fmt.Printf("no embedded parts\r\n")
		t.Fail()
	} else {
		embeddedPart1 := importedEmbeddedParts[0]
		importedEmbedded := embeddedPart1.GetBody()
		if embedded != importedEmbedded {
			fmt.Printf("`embedded` corrupted. expected \"%s\", got \"%s\"\r\n", embedded, importedEmbedded)
			t.Fail()
		}
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

func TestRoundTripTextAndHtmlAndEmbeddedAndAttachment(t *testing.T) {
	email := emptyEmail.
		WithText(text).
		WithHtml(html).
		WithEmbeddedString(embedded).
		WithAttachedString(attached)

	emailString := email.String()

	importedEmail, err := simpleemail.Import(emailString)
	if err != nil {
		fmt.Println(`import failed:`)
		fmt.Println(err)
		t.FailNow()
	}

	importedText := importedEmail.GetText()
	if importedText != text {
		fmt.Printf("`text` corrupted. expected \"%s\", got \"%s\"\r\n", text, importedText)
		t.Fail()
	}

	importedHtml := importedEmail.GetHtml()
	if importedHtml != html {
		fmt.Printf("`html` corrupted. expected \"%s\", got \"%s\"\r\n", html, importedHtml)
		t.Fail()
	}

	importedEmbeddedParts := importedEmail.GetEmbedded()
	if len(importedEmbeddedParts) < 1 {
		fmt.Printf("no embedded parts\r\n")
		t.Fail()
	} else {
		embeddedPart1 := importedEmbeddedParts[0]
		importedEmbedded := embeddedPart1.GetBody()
		if embedded != importedEmbedded {
			fmt.Printf("`embedded` corrupted. expected \"%s\", got \"%s\"\r\n", embedded, importedEmbedded)
			t.Fail()
		}
	}

	importedAttachedParts := importedEmail.GetAttachments()
	if len(importedAttachedParts) < 1 {
		fmt.Printf("no attached parts\r\n")
		t.Fail()
	} else {
		attachedPart1 := importedAttachedParts[0]
		importAttached := attachedPart1.GetBody()
		if attached != importAttached {
			fmt.Printf("`attached` corrupted. expected \"%s\", got \"%s\"\r\n", attached, importAttached)
			t.Fail()
		}
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
	//email.GetSubject()
}

func TestRoundTripFull(t *testing.T) {

	email := emptyEmail.
		WithText(text).
		WithHtml(html).
		WithEmbeddedString(embedded).
		WithAttachedString(attached)

	emailString := email.String()

	importedEmail, err := simpleemail.Import(emailString)
	if err != nil {
		fmt.Println(`import failed:`)
		fmt.Println(err)
		t.FailNow()
	}

	importedFrom := importedEmail.GetFrom()
	if !addressSlicesEqual(from, importedFrom) {
		fmt.Printf("`From` corrupted. expected %#v, got %#v\r\n", from, importedFrom)
		t.Fail()
	}

	importedTo := importedEmail.GetTo()
	if !addressSlicesEqual(to, importedTo) {
		fmt.Printf("`To` corrupted. expected %#v, got %#v\r\n", to, importedTo)
		t.Fail()
	}

	importedCc := importedEmail.GetCc()
	if !addressSlicesEqual(cc, importedCc) {
		fmt.Printf("`Cc` corrupted. expected %#v, got %#v\r\n", cc, importedCc)
		t.Fail()
	}

	importedBcc := importedEmail.GetBcc()
	if !addressSlicesEqual(bcc, importedBcc) {
		fmt.Printf("`Bcc` corrupted. expected %#v, got %#v\r\n", bcc, importedBcc)
		t.Fail()
	}

	importedSubject := importedEmail.GetSubject()
	if importedSubject != subject {
		fmt.Printf("`Subject` corrupted. expected \"%s\", got \"%s\"\r\n", subject, importedSubject)
		t.Fail()
	}

	importedText := importedEmail.GetText()
	if importedText != text {
		fmt.Printf("`text` corrupted. expected \"%s\", got \"%s\"\r\n", text, importedText)
		t.Fail()
	}

	importedHtml := importedEmail.GetHtml()
	if importedHtml != html {
		fmt.Printf("`html` corrupted. expected \"%s\", got \"%s\"\r\n", html, importedHtml)
		t.Fail()
	}

	importedEmbeddedParts := importedEmail.GetEmbedded()
	if len(importedEmbeddedParts) < 1 {
		fmt.Printf("no embedded parts\r\n")
		t.Fail()
	} else {
		embeddedPart1 := importedEmbeddedParts[0]
		importedEmbedded := embeddedPart1.GetBody()
		if embedded != importedEmbedded {
			fmt.Printf("`embedded` corrupted. expected \"%s\", got \"%s\"\r\n", embedded, importedEmbedded)
			t.Fail()
		}
	}

	importedAttachedParts := importedEmail.GetAttachments()
	if len(importedAttachedParts) < 1 {
		fmt.Printf("no attached parts\r\n")
		t.Fail()
	} else {
		attachedPart1 := importedAttachedParts[0]
		importAttached := attachedPart1.GetBody()
		if attached != importAttached {
			fmt.Printf("`attached` corrupted. expected \"%s\", got \"%s\"\r\n", attached, importAttached)
			t.Fail()
		}
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
	//email.GetSubject()
}
