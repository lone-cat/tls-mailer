package simpleemail

import (
	"github.com/lone-cat/tls-mailer/simpleemail/test"
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
	for _, em := range test.emailsForTest {
		testCompare(em, t)
	}
}

func testCompare(email *Email, t *testing.T) {
	emailString := email.String()

	importedEmail, err := Import(emailString)
	if err != nil {
		t.Fatalf("%s\r\n%#v", `import failed:`, err)
	}

	errors := emailsDiffErrors(email, importedEmail)
	for _, erro := range errors {
		t.Error(erro)
	}
}
