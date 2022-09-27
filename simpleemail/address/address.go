package address

import (
	"fmt"
	"net/mail"
)

type Address interface {
	fmt.Stringer
	GetName() string
	WithName(string) Address
	GetEmail() string
	WithEmail(string) Address
	ToMailAddress() *mail.Address
}

type addr struct {
	name  string
	email string
}

func NewAddress(name string, email string) Address {
	return &addr{name: name, email: email}
}

func NewAddressFromMailAddr(addrs *mail.Address) Address {
	return &addr{name: addrs.Name, email: addrs.Address}
}

func (a *addr) GetName() string {
	return a.name
}

func (a *addr) WithName(name string) Address {
	return &addr{
		name:  name,
		email: a.email,
	}
}

func (a *addr) GetEmail() string {
	return a.email
}

func (a *addr) WithEmail(email string) Address {
	return &addr{
		name:  a.name,
		email: email,
	}
}

func (a *addr) ToMailAddress() *mail.Address {
	return &mail.Address{
		Address: a.email,
		Name:    a.name,
	}
}

func (a *addr) String() string {
	return a.ToMailAddress().String()
}
