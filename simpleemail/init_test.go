package simpleemail_test

import (
	"fmt"
	"github.com/lone-cat/stackerrors"
)

func init() {
	stackerrors.SetDebugMode(true)
	fmt.Println(`debug mode`)
}
