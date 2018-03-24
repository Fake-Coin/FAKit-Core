package FAKitCore

// #cgo CFLAGS: -I./secp256k1 -I./secp256k1/include -DBITCOIN_TEST_NO_MAIN
// #include <stdlib.h>
// #include <BRBase58.h>
import "C"
import "unsafe"

func BRBase58Encode(data []byte) string {
	l := C.BRBase58Encode(nil, 0, (*C.uint8_t)(&data[0]), C.size_t(len(data)))
	str := make([]byte, int(l))
	l = C.BRBase58Encode((*C.char)(unsafe.Pointer(&str[0])), l, (*C.uint8_t)(&data[0]), C.size_t(len(data)))
	return string(str[:int(l)])
}

func BRBase58Decode(str string) []byte {
	str_c := C.CString(str)
	defer C.free(unsafe.Pointer(str_c))

	l := C.BRBase58Decode(nil, 0, str_c)
	data := make([]byte, int(l))
	l = C.BRBase58Decode((*C.uint8_t)(&data[0]), l, str_c)
	return data[:int(l)]
}

func BRBase58CheckEncode(data []byte) string {
	l := C.BRBase58CheckEncode(nil, 0, (*C.uint8_t)(&data[0]), C.size_t(len(data)))
	str := make([]byte, int(l))
	l = C.BRBase58CheckEncode((*C.char)(unsafe.Pointer(&str[0])), l, (*C.uint8_t)(&data[0]), C.size_t(len(data)))
	return string(str[:int(l)])
}

func BRBase58CheckDecode(str string) []byte {
	str_c := C.CString(str)
	defer C.free(unsafe.Pointer(str_c))

	l := C.BRBase58CheckDecode(nil, 0, str_c)
	data := make([]byte, int(l))
	l = C.BRBase58CheckDecode((*C.uint8_t)(&data[0]), l, str_c)
	return data[:int(l)]
}
