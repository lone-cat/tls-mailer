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

func TestRoundTrip(t *testing.T) {
	addr1 := mail.Address{Name: `Иванов Иван`, Address: `first@email.addr`}
	addr2 := mail.Address{Name: `Петров Петр`, Address: `second@email.addr`}
	addr3 := mail.Address{Name: `Сидоров Сидор`, Address: `third@email.addr`}
	addr4 := mail.Address{Name: `Иванов Петр`, Address: `fourth@email.addr`}
	addr5 := mail.Address{Name: `Иванов Сидор`, Address: `fifth@email.addr`}
	addr6 := mail.Address{Name: `Петров Иван`, Address: `sixth@email.addr`}
	addr7 := mail.Address{Name: `Петров Сидор`, Address: `seventh@email.addr`}
	addr8 := mail.Address{Name: `Сидоров Иван`, Address: `eighth@email.addr`}
	addr9 := mail.Address{Name: `Сидоров Петр`, Address: `ninth@email.addr`}

	from := []mail.Address{addr1, addr2}
	to := []mail.Address{addr2, addr3, addr4, addr5}
	cc := []mail.Address{addr6, addr7, addr8}
	bcc := []mail.Address{addr9, addr1}

	subject := "Какая то странная длинная тема s angliiskimi simvolami так чтобы \r\nточно поместилась только на несколько строк"

	text := `Какой то не менее длинный текст, чтобы он тоже был на несколько строк, но при этом еще длиннее чем предыдущий` + "\r\n" +
		`Кстати, этот текст еще и будет иметь перевод строки. Tak zhe on soderzhit английские буквы )`
	html := `<h1>Здесь длина текста уже не будет иметь значения</h1>`

	embedded := `aaa`

	attached := `bbb`

	email := simpleemail.NewEmptyEmail().
		WithFrom(from).
		WithTo(to).
		WithCc(cc).
		WithBcc(bcc).
		WithSubject(subject).
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
