package FAKitCore

// #cgo CFLAGS: -I./secp256k1 -I./secp256k1/include -DBITCOIN_TEST_NO_MAIN
// #include <stdlib.h>
// #include <BRAddress.h>
import "C"
import "unsafe"

func BRAddressIsValid(addr string) bool {
	addr_c := C.CString(addr)
	defer C.free(unsafe.Pointer(addr_c))
	return C.BRAddressIsValid(addr_c) == 1
}
