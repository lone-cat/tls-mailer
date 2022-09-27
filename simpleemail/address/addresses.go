package address

import (
	"fmt"
	"github.com/lone-cat/tls-mailer/common"
	"net/mail"
	"strings"
)

type AddressList interface {
	fmt.Stringer
	ExportAddressSlice() []Address
	ExportMailAddressSlice() []*mail.Address
	WithAddressSlice(...Address) AddressList
	WithMailAddressSlice(...*mail.Address) AddressList
}

type addressList struct {
	addrs []Address
}

func NewAddressList(addrs ...Address) AddressList {
	newAddressList := &addressList{
		addrs: common.CloneSlice(addrs),
	}

	return newAddressList
}

func NewAddressListFromMailAddresses(addrs ...*mail.Address) AddressList {
	return NewAddressList(mailAddrSliceToAddrSlice(addrs...)...)
}

func (a *addressList) ExportAddressSlice() []Address {
	return common.CloneSlice(a.addrs)
}

func (a *addressList) ExportMailAddressSlice() []*mail.Address {
	return addrSliceToMailAddrSlice(a.addrs...)
}

func (a *addressList) WithAddressSlice(addrs ...Address) AddressList {
	return NewAddressList(addrs...)
}

func (a *addressList) WithMailAddressSlice(addrs ...*mail.Address) AddressList {
	return NewAddressListFromMailAddresses(addrs...)
}

func (a *addressList) String() string {
	strAddr := make([]string, len(a.addrs))
	for i := range a.addrs {
		strAddr[i] = a.addrs[i].String()
	}

	return strings.Join(strAddr, `, `)
}

func mailAddrSliceToAddrSlice(mailAddresses ...*mail.Address) []Address {
	addresses := make([]Address, len(mailAddresses))
	for i := range mailAddresses {
		addresses[i] = NewAddressFromMailAddr(mailAddresses[i])
	}

	return addresses
}

func addrSliceToMailAddrSlice(addresses ...Address) []*mail.Address {
	mailAddresses := make([]*mail.Address, len(addresses))
	for i := range mailAddresses {
		mailAddresses[i] = addresses[i].ToMailAddress()
	}

	return mailAddresses
}
