package FAKitCore

// #cgo CFLAGS: -I./secp256k1 -I./secp256k1/include -DBITCOIN_TEST_NO_MAIN
// #include <stdlib.h>
// #include <BRAddress.h>
// #include <BRBIP32Sequence.h>
// #include <BRBIP39Mnemonic.h>
import "C"
import "unsafe"

type (
	PrivateKey [64]uint8
	PublicKey  [33]uint8
	UInt256    [32]uint8
)

func BRBIP39DeriveKey(keyPhrase string) (k PrivateKey) {
	phrase_c := C.CString(keyPhrase)
	defer C.free(unsafe.Pointer(phrase_c))

	C.BRBIP39DeriveKey(unsafe.Pointer(&k[0]), phrase_c, nil)
	return k
}

type BRMasterPubKey C.BRMasterPubKey

func (mpk BRMasterPubKey) FingerPrint() uint32 {
	return uint32(mpk.fingerPrint)
}

func (mpk BRMasterPubKey) ChainCode() UInt256 {
	return UInt256(mpk.chainCode)
}

func (mpk BRMasterPubKey) PubKey() PublicKey {
	return *(*PublicKey)(unsafe.Pointer(&mpk.pubKey[0]))
}

func BRBIP32MasterPubKey(priv PrivateKey) BRMasterPubKey {
	return BRMasterPubKey(C.BRBIP32MasterPubKey(unsafe.Pointer(&priv[0]), C.size_t(len(priv))))
}
