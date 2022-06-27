package simpleemail_test

import (
	"fmt"
	"github.com/lone-cat/tls-mailer/simpleemail"
	"net/mail"
	"testing"
)

func addressSlicesEqual(addressSlice1 []*mail.Address, addressSlice2 []*mail.Address) bool {
	addressMap1 := convertAddressSliceToMapByEmail(addressSlice1)
	addressMap2 := convertAddressSliceToMapByEmail(addressSlice2)
	if len(addressMap1) != len(addressMap2) {
		return false
	}

	var ok bool
	var addr2 *mail.Address
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

func convertAddressSliceToMapByEmail(addressSlice []*mail.Address) map[string]*mail.Address {
	result := make(map[string]*mail.Address)
	for _, addr := range addressSlice {
		result[addr.Address] = &mail.Address{Name: addr.Name, Address: addr.Name}
	}
	return result
}

func TestRoundTrip(t *testing.T) {
	for _, em := range emailsForTest {
		testCompare(em, t)
	}
}

func testCompare(email *simpleemail.Email, t *testing.T) {
	emailString := email.String()

	importedEmail, err := simpleemail.Import(emailString)
	if err != nil {
		fmt.Println(`import failed:`)
		fmt.Println(err)
		t.FailNow()
	}

	if !addressSlicesEqual(email.GetFrom(), importedEmail.GetFrom()) {
		fmt.Printf("`From` does not match: expected %#v, got %#v\r\n", email.GetFrom(), importedEmail.GetFrom())
		t.Fail()
	}
	if !addressSlicesEqual(email.GetTo(), importedEmail.GetTo()) {
		fmt.Printf("`To` does not match: expected %#v, got %#v\r\n", email.GetTo(), importedEmail.GetTo())
		t.Fail()
	}
	if !addressSlicesEqual(email.GetCc(), importedEmail.GetCc()) {
		fmt.Printf("`Cc` does not match: expected %#v, got %#v\r\n", email.GetCc(), importedEmail.GetCc())
		t.Fail()
	}
	if !addressSlicesEqual(email.GetBcc(), importedEmail.GetBcc()) {
		fmt.Printf("`Bcc` does not match: expected %#v, got %#v\r\n", email.GetBcc(), importedEmail.GetBcc())
		t.Fail()
	}
	if email.GetSubject() != importedEmail.GetSubject() {
		fmt.Printf("`subject` does not match: expected `%s`, got `%s`\r\n", email.GetSubject(), importedEmail.GetSubject())
		t.Fail()
	}
	if email.GetText() != importedEmail.GetText() {
		fmt.Printf("`text` does not match: expected `%s`, got `%s`\r\n", email.GetText(), importedEmail.GetText())
		t.Fail()
	}
	if email.GetHtml() != importedEmail.GetHtml() {
		fmt.Printf("`html` does not match: expected `%s`, got `%s`\r\n", email.GetHtml(), importedEmail.GetHtml())
		t.Fail()
	}
	if len(email.GetEmbedded()) != len(importedEmail.GetEmbedded()) {
		fmt.Printf("`Embedded` count does not match: expected `%d`, got `%d`\r\n", len(email.GetEmbedded()), len(importedEmail.GetEmbedded()))
		t.Fail()
	} else {
		if len(email.GetEmbedded()) > 0 {
			if email.GetEmbedded()[0].GetBody() != importedEmail.GetEmbedded()[0].GetBody() {
				fmt.Printf("`Embedded` body does not match: expected `%s`, got `%s`\r\n", email.GetEmbedded()[0].GetBody(), importedEmail.GetEmbedded()[0].GetBody())
				t.Fail()
			}
		}
	}
	if len(email.GetAttachments()) != len(importedEmail.GetAttachments()) {
		fmt.Printf("`Attachments` count does not match: expected `%d`, got `%d`\r\n", len(email.GetAttachments()), len(importedEmail.GetAttachments()))
		t.Fail()
	} else {
		if len(email.GetAttachments()) > 0 {
			if email.GetAttachments()[0].GetBody() != importedEmail.GetAttachments()[0].GetBody() {
				fmt.Printf("`Attachments` body does not match: expected `%s`, got `%s`\r\n", email.GetAttachments()[0].GetBody(), importedEmail.GetAttachments()[0].GetBody())
				t.Fail()
			}
		}
	}
}
