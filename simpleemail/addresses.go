package simpleemail

import "net/mail"

type addresses []mail.Address

func newAddresses() addresses {
	return make([]mail.Address, 0)
}

func (a addresses) clone() addresses {
	clonedAddresses := make([]mail.Address, len(a))
	for index, address := range a {
		clonedAddresses[index] = mail.Address{Name: address.Name, Address: address.Address}
	}
	return clonedAddresses
}
