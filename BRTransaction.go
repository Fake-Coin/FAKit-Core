package FAKitCore

// #cgo CFLAGS: -I./secp256k1 -I./secp256k1/include -DBITCOIN_TEST_NO_MAIN
// #include <stdlib.h>
// #include <BRTransaction.h>
import "C"
import (
	"time"
	"unsafe"
)

type BRTransaction C.BRTransaction

func BRTransactionNew() *BRTransaction {
	return (*BRTransaction)(C.BRTransactionNew())
}

func BRTransactionParse(txSerial []byte) *BRTransaction {
	return (*BRTransaction)(C.BRTransactionParse((*C.uint8_t)(unsafe.Pointer(&txSerial[0])), C.size_t(len(txSerial))))
}

func (tx *BRTransaction) Version() uint32 {
	return uint32(tx.version)
}

func (tx *BRTransaction) BlockHeight() uint32 {
	return uint32(tx.blockHeight)
}

func (tx *BRTransaction) Time() time.Time {
	return time.Unix(int64(tx.timestamp), 0)
}

func (tx *BRTransaction) Hash() UInt256 {
	return UInt256(tx.txHash)
}

func (tx *BRTransaction) Serialize() []byte {
	l := C.BRTransactionSerialize((*C.BRTransaction)(tx), nil, 0)
	data := make([]byte, int(l))
	l = C.BRTransactionSerialize((*C.BRTransaction)(tx), (*C.uint8_t)(unsafe.Pointer(&data[0])), l)
	return data[:int(l)]
}

func (tx *BRTransaction) ShuffleOutputs() {
	C.BRTransactionShuffleOutputs((*C.BRTransaction)(tx))
}

func (tx *BRTransaction) Size() int {
	return int(C.BRTransactionSize((*C.BRTransaction)(tx)))
}

func (tx *BRTransaction) StandardFee() uint64 {
	return uint64(C.BRTransactionStandardFee((*C.BRTransaction)(tx)))
}

func (tx *BRTransaction) IsSigned() bool {
	return C.BRTransactionIsSigned((*C.BRTransaction)(tx)) == 1
}

func (tx *BRTransaction) DeepCopy() *BRTransaction {
	return (*BRTransaction)(C.BRTransactionCopy((*C.BRTransaction)(tx)))
}

func (tx *BRTransaction) Free() {
	C.BRTransactionFree((*C.BRTransaction)(tx))
}
