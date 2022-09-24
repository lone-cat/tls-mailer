package addresses

import "net/mail"

type addresses struct {
	addrs []*mail.Address
}

func NewAddressList(addrs ...*mail.Address) AddressList {
	newAddressList := &addresses{
		addrs: duplicateSlice(addrs),
	}

	return newAddressList
}

func (a *addresses) ExportAddressSlice() []*mail.Address {
	return duplicateSlice(a.addrs)
}

func (a *addresses) WithAddressSlice(addrs []*mail.Address) AddressList {
	return NewAddressList(addrs...)
}

func duplicateSlice(src []*mail.Address) (dst []*mail.Address) {
	dst = make([]*mail.Address, len(src))
	copy(dst, src)

	return
}
