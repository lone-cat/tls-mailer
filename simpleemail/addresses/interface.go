package addresses

import "net/mail"

type AddressList interface {
	ExportAddressSlice() []*mail.Address
	WithAddressSlice(addrs []*mail.Address) AddressList
}
