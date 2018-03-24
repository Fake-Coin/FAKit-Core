package FAKitCore

// #cgo CFLAGS: -I./secp256k1 -I./secp256k1/include -DBITCOIN_TEST_NO_MAIN
// #include <stdlib.h>
// #include <stdio.h>
// #include <BRPeerManager.h>
//
// extern void syncStarted(void *info);
// extern void syncStopped(void *info, int error);
// extern void txStatusUpdate(void *info);
// extern void saveBlocks(void *info, int replace, BRMerkleBlock **blocks, size_t blocksCount);
// extern void savePeers(void *info, int replace, BRPeer *peers, size_t peersCount);
// extern int networkIsReachable(void *info);
// extern void threadCleanup(void *info);
import "C"
import (
	"unsafe"
)

type BRPeerManager struct {
	brmgr *C.BRPeerManager

	SyncStarted        func()
	SyncStopped        func(errCode int)
	TxStatusUpdate     func()
	SaveBlocks         func(replace bool, blocks []*BRMerkleBlock)
	SavePeers          func(replace bool, peers []BRPeer)
	NetworkIsReachable func() int
	ThreadCleanup      func()
}

type BRMerkleBlock = *C.BRMerkleBlock
type BRPeer C.BRPeer

func BRPeerManagerNewMainNet(w *BRWallet, keyTime uint32, blocks []BRMerkleBlock, peers []BRPeer) *BRPeerManager {
	var cblocks *BRMerkleBlock
	if 0 < len(blocks) {
		cblocks = &blocks[0]
	}

	var cpeers *BRPeer
	if 0 < len(peers) {
		cpeers = &peers[0]
	}

	mgr := &BRPeerManager{
		brmgr: C.BRPeerManagerNewMainNet(w.brwallet, C.uint32_t(keyTime),
			cblocks, C.size_t(len(blocks)),
			(*C.BRPeer)(unsafe.Pointer(cpeers)), C.size_t(len(peers))),
	}

	C.BRPeerManagerSetCallbacks(mgr.brmgr, unsafe.Pointer(mgr),
		(*[0]byte)(C.syncStarted), (*[0]byte)(C.syncStopped), (*[0]byte)(C.txStatusUpdate),
		(*[0]byte)(C.saveBlocks), (*[0]byte)(C.savePeers),
		(*[0]byte)(C.networkIsReachable), (*[0]byte)(C.threadCleanup))

	return mgr
}

func (w *BRPeerManager) Connect() {
	C.BRPeerManagerConnect(w.brmgr)
}

func (w *BRPeerManager) Disconnect() {
	C.BRPeerManagerDisconnect(w.brmgr)
}

func (w *BRPeerManager) Rescan() {
	C.BRPeerManagerRescan(w.brmgr)
}

func (w *BRPeerManager) Status() BRPeerStatus {
	return C.BRPeerManagerConnectStatus(w.brmgr)
}

func (w *BRPeerManager) PublishTx(tx *BRTransaction) {
	C.BRPeerManagerPublishTx(w.brmgr, (*C.BRTransaction)(tx), nil, nil)
}

func (w *BRPeerManager) EstimatedBlockHeight() uint32 {
	return uint32(C.BRPeerManagerEstimatedBlockHeight(w.brmgr))
}

func (w *BRPeerManager) LastBlockHeight() uint32 {
	return uint32(C.BRPeerManagerLastBlockHeight(w.brmgr))
}

func (w *BRPeerManager) LastBlockTimestamp() uint32 {
	return uint32(C.BRPeerManagerLastBlockTimestamp(w.brmgr))
}

func (w *BRPeerManager) Progress() float64 {
	return float64(C.BRPeerManagerSyncProgress(w.brmgr, 0))
}

func (w *BRPeerManager) PeerCount() int {
	return int(C.BRPeerManagerPeerCount(w.brmgr))
}

func (w *BRPeerManager) PeerName() string {
	return C.GoString(C.BRPeerManagerDownloadPeerName(w.brmgr))
}

func (w *BRPeerManager) Free() {
	C.BRPeerManagerFree(w.brmgr)
}

type BRPeerStatus = C.BRPeerStatus

const (
	StatusDisconnected = C.BRPeerStatusDisconnected
	StatusConnecting   = C.BRPeerStatusConnecting
	StatusConnected    = C.BRPeerStatusConnected
)

//export syncStarted
func syncStarted(info unsafe.Pointer) {
	if mgr := (*BRPeerManager)(info); mgr.SyncStarted != nil {
		mgr.SyncStarted()
	}
}

//export syncStopped
func syncStopped(info unsafe.Pointer, errCode C.int) {
	if mgr := (*BRPeerManager)(info); mgr.SyncStopped != nil {
		mgr.SyncStopped(int(errCode))
	}
}

//export txStatusUpdate
func txStatusUpdate(info unsafe.Pointer) {
	if mgr := (*BRPeerManager)(info); mgr.TxStatusUpdate != nil {
		mgr.TxStatusUpdate()
	}
}

//export saveBlocks
func saveBlocks(info unsafe.Pointer, replace C.int, blocks **C.BRMerkleBlock, blocksCount C.size_t) {
	if mgr := (*BRPeerManager)(info); mgr.SaveBlocks != nil {
		mgr.SaveBlocks(replace == 1, nil)
	}
}

//export savePeers
func savePeers(info unsafe.Pointer, replace C.int, peers *C.BRPeer, peersCount C.size_t) {
	if mgr := (*BRPeerManager)(info); mgr.SavePeers != nil {
		mgr.SavePeers(replace == 1, nil)
	}
}

//export networkIsReachable
func networkIsReachable(info unsafe.Pointer) C.int {
	if mgr := (*BRPeerManager)(info); mgr.NetworkIsReachable != nil {
		return C.int(mgr.NetworkIsReachable())
	}
	return 1
}

//export threadCleanup
func threadCleanup(info unsafe.Pointer) {
	if mgr := (*BRPeerManager)(info); mgr.ThreadCleanup != nil {
		mgr.ThreadCleanup()
	}
}
