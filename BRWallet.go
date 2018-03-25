package FAKitCore

// #cgo CFLAGS: -I./secp256k1 -I./secp256k1/include -DBITCOIN_TEST_NO_MAIN
// #include <stdlib.h>
// #include <BRWallet.h>
//
// extern void balanceChanged(void *info, uint64_t balance);
// extern void txAdded(void *info, BRTransaction *tx);
// extern void txUpdated(void *info, UInt256 txHashes[], size_t txCount, uint32_t blockHeight, uint32_t timestamp);
// extern void txDeleted(void *info, UInt256 txHash, int notifyUser, int recommendRescan);
import "C"
import (
	"unsafe"
)

// BRWallet *BRWalletNew(BRTransaction *transactions[], size_t txCount, BRMasterPubKey mpk);

type BRWallet struct {
	brwallet *C.BRWallet

	BalanceChanged func(balance uint64)
	TXAdded        func(tx *BRTransaction)
	TXUpdated      func(txHashes []UInt256, blockHeight uint32, timestamp uint32)
	TXDeleted      func(txHash UInt256, notifyUser bool, recommendRescan bool)
}

func BRWalletNew(txs []*BRTransaction, mpk BRMasterPubKey) (w *BRWallet) {
	var ctxs **C.BRTransaction
	if 0 < len(txs) {
		ctxs = (**C.BRTransaction)(unsafe.Pointer(&txs[0]))
	}

	wallet := &BRWallet{brwallet: C.BRWalletNew(ctxs, C.size_t(len(txs)), C.BRMasterPubKey(mpk))}
	C.BRWalletSetCallbacks(wallet.brwallet, unsafe.Pointer(wallet),
		(*[0]byte)(C.balanceChanged), (*[0]byte)(C.txAdded),
		(*[0]byte)(C.txUpdated), (*[0]byte)(C.txDeleted))

	return wallet
} // result must be freed using BRTransactionFree()

func (w *BRWallet) ReceiveAddress() string {
	addr := C.BRWalletReceiveAddress(w.brwallet)
	return C.GoString(&addr.s[0])
}

func (w *BRWallet) Balance() uint64 {
	return uint64(C.BRWalletBalance(w.brwallet))
}

func (w *BRWallet) CreateTransaction(amount uint64, addr string) *BRTransaction {
	addr_c := C.CString(addr)
	defer C.free(unsafe.Pointer(addr_c))

	return (*BRTransaction)(C.BRWalletCreateTransaction(w.brwallet, C.uint64_t(amount), addr_c))
}

// int BRWalletSignTransaction(
// 	BRWallet *wallet,
// 	BRTransaction *tx,
// 	int forkId,
// 	const void *seed,
// 	size_t seedLen);
func (w *BRWallet) SignTransaction(tx *BRTransaction, priv PrivateKey) bool {
	return C.BRWalletSignTransaction(w.brwallet, (*C.BRTransaction)(tx), 0x00,
		unsafe.Pointer(&priv[0]),
		C.size_t(len(priv))) == 1
}

func (w *BRWallet) TotalSent() uint64 {
	return uint64(C.BRWalletTotalSent(w.brwallet))
}

func (w *BRWallet) TotalReceived() uint64 {
	return uint64(C.BRWalletTotalReceived(w.brwallet))
}

func (w *BRWallet) SetFeePerKB(n uint64) {
	C.BRWalletSetFeePerKb(w.brwallet, C.uint64_t(n))
}

func (w *BRWallet) FeeForTxSize(b int) uint64 {
	return uint64(C.BRWalletFeeForTxSize(w.brwallet, C.size_t(b)))
}

func (w *BRWallet) FeeForTxAmount(amount uint64) uint64 {
	return uint64(C.BRWalletFeeForTxAmount(w.brwallet, C.uint64_t(amount)))
}

func (w *BRWallet) MinOutput() uint64 {
	return uint64(C.BRWalletMinOutputAmount(w.brwallet))
}

func (w *BRWallet) MaxOutput() uint64 {
	return uint64(C.BRWalletMaxOutputAmount(w.brwallet))
}

func (w *BRWallet) Free() {
	C.BRWalletFree(w.brwallet)
}

//export balanceChanged
func balanceChanged(info unsafe.Pointer, balance C.uint64_t) {
	if wallet := (*BRWallet)(info); wallet.BalanceChanged != nil {
		wallet.BalanceChanged(uint64(balance))
	}
}

//export txAdded
func txAdded(info unsafe.Pointer, tx *BRTransaction) {
	if wallet := (*BRWallet)(info); wallet.TXAdded != nil {
		wallet.TXAdded(tx)
	}
}

//export txUpdated
func txUpdated(info unsafe.Pointer, txHashes *C.UInt256, txCount C.size_t, blockHeight C.uint32_t, timestamp C.uint32_t) {
	if wallet := (*BRWallet)(info); wallet.TXAdded != nil {
		cHash := (*[1 << 30]UInt256)(unsafe.Pointer(txHashes))[:int(txCount):int(txCount)]

		// copy from c stack to go heap
		goHash := make([]UInt256, len(cHash))
		for i, h := range cHash {
			goHash[i] = h
		}

		wallet.TXUpdated(goHash, uint32(blockHeight), uint32(timestamp))
	}
}

//export txDeleted
func txDeleted(info unsafe.Pointer, txHash C.UInt256, notifyUser C.int, recommendRescan C.int) {
	if wallet := (*BRWallet)(info); wallet.TXDeleted != nil {
		wallet.TXDeleted(UInt256(txHash), notifyUser == 1, recommendRescan == 1)
	}
}
