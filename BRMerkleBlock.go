package FAKitCore

// #cgo CFLAGS: -I./secp256k1 -I./secp256k1/include -DBITCOIN_TEST_NO_MAIN
// #include <stdlib.h>
// #include <BRMerkleBlock.h>
import "C"
import (
	"time"
	"unsafe"
)

type BRMerkleBlock C.BRMerkleBlock

func BRMerkleBlockNew() *BRMerkleBlock {
	return (*BRMerkleBlock)(C.BRMerkleBlockNew())
}

func BRMerkleBlockParse(blkSerial []byte) *BRMerkleBlock {
	return (*BRMerkleBlock)(C.BRMerkleBlockParse((*C.uint8_t)(unsafe.Pointer(&blkSerial[0])), C.size_t(len(blkSerial))))
}

func (blk *BRMerkleBlock) TxHashes() []UInt256 {
	l := C.BRMerkleBlockTxHashes((*C.BRMerkleBlock)(blk), nil, 0)
	txHashes := make([]UInt256, int(l))
	l = C.BRMerkleBlockTxHashes((*C.BRMerkleBlock)(blk), (*C.UInt256)(unsafe.Pointer(&txHashes[0])), l)
	return txHashes[:int(l)]
}

func (blk *BRMerkleBlock) SetTxHashes(hashes []UInt256, flags []uint8) {
	C.BRMerkleBlockSetTxHashes(
		(*C.BRMerkleBlock)(blk),
		(*C.UInt256)(unsafe.Pointer(&hashes[0])), C.size_t(len(hashes)),
		(*C.uint8_t)(unsafe.Pointer(&flags[0])), C.size_t(len(flags)))
}

func (blk *BRMerkleBlock) IsValid() bool {
	t := time.Now().Unix()
	return C.BRMerkleBlockIsValid((*C.BRMerkleBlock)(blk), C.uint32_t(t)) == 1
}

func (blk *BRMerkleBlock) VerifyDifficulty(previous *BRMerkleBlock, transitionTime uint32) bool {
	return C.BRMerkleBlockVerifyDifficulty((*C.BRMerkleBlock)(blk), (*C.BRMerkleBlock)(previous), C.uint32_t(transitionTime)) == 1
}

func (blk *BRMerkleBlock) ContainsTxHash(txHash UInt256) bool {
	return C.BRMerkleBlockContainsTxHash((*C.BRMerkleBlock)(blk), C.UInt256(txHash)) == 1
}

func (blk *BRMerkleBlock) DeepCopy() *BRMerkleBlock {
	return (*BRMerkleBlock)(C.BRMerkleBlockCopy((*C.BRMerkleBlock)(blk)))
}

func (blk *BRMerkleBlock) Serialize() []byte {
	l := C.BRMerkleBlockSerialize((*C.BRMerkleBlock)(blk), nil, 0)
	data := make([]byte, int(l))
	l = C.BRMerkleBlockSerialize((*C.BRMerkleBlock)(blk), (*C.uint8_t)(unsafe.Pointer(&data[0])), l)
	return data[:int(l)]
}

func (blk *BRMerkleBlock) Free() {
	C.BRMerkleBlockFree((*C.BRMerkleBlock)(blk))
}
