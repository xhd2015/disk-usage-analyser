//go:build darwin

package analyse

import (
	"encoding/binary"
	"syscall"
	"unsafe"
)

const (
	attrBitMapCount        = 5
	fsoptAttrCmnExtended   = 0x00000020
	attrCmnDocumentID      = 0x00100000
	attrCmnextCloneID      = 0x00000100
	attrCmnextCloneRefCnt  = 0x00001000
	sysGetattrlist         = 220
)

type attrList struct {
	bitmapcount uint16
	reserved    uint16
	commonattr  uint32
	volattr     uint32
	dirattr     uint32
	fileattr    uint32
	forkattr    uint32
}

type cloneGroupKey struct {
	docID   uint32
	cloneID uint64
}

type cloneGroupState struct {
	inodes map[inodeKey]struct{}
	size   int64
}

type cloneGroupTracker struct {
	groups         map[cloneGroupKey]*cloneGroupState
	countedForSize map[cloneGroupKey]struct{}
}

func newCloneGroupTracker() *cloneGroupTracker {
	return &cloneGroupTracker{
		groups:         make(map[cloneGroupKey]*cloneGroupState),
		countedForSize: make(map[cloneGroupKey]struct{}),
	}
}

func (t *cloneGroupTracker) CountSize(path string, size int64) int64 {
	docID, cloneID, _, ok := fileCloneAttrs(path)
	if !ok {
		return size
	}
	key, ok := cloneGroupKeyFrom(docID, cloneID)
	if !ok {
		return size
	}
	if _, counted := t.countedForSize[key]; counted {
		return 0
	}
	t.countedForSize[key] = struct{}{}
	return size
}

func (t *cloneGroupTracker) Add(path string, inode inodeKey, size int64) {
	docID, cloneID, _, ok := fileCloneAttrs(path)
	if !ok {
		return
	}

	key, ok := cloneGroupKeyFrom(docID, cloneID)
	if !ok {
		return
	}

	group := t.groups[key]
	if group == nil {
		group = &cloneGroupState{inodes: make(map[inodeKey]struct{})}
		t.groups[key] = group
	}
	group.inodes[inode] = struct{}{}
	if group.size == 0 {
		group.size = size
	}
}

func (t *cloneGroupTracker) TotalSharedCloneSize() int64 {
	var total int64
	for _, group := range t.groups {
		localRefs := len(group.inodes)
		if localRefs <= 1 {
			continue
		}
		// Subtree-local refs only. APFS global clone ref counts include store
		// and other projects outside the analysed path and must not be used here.
		total += group.size * int64(localRefs)
	}
	return total
}

func cloneGroupKeyFrom(docID uint32, cloneID uint64) (cloneGroupKey, bool) {
	if docID != 0 {
		return cloneGroupKey{docID: docID}, true
	}
	if cloneID != 0 {
		return cloneGroupKey{cloneID: cloneID}, true
	}
	return cloneGroupKey{}, false
}

func fileCloneAttrs(path string) (docID uint32, cloneID uint64, refCount uint32, ok bool) {
	pathPtr, err := syscall.BytePtrFromString(path)
	if err != nil {
		return 0, 0, 0, false
	}

	al := attrList{
		bitmapcount: attrBitMapCount,
		commonattr:  attrCmnDocumentID,
		forkattr:    attrCmnextCloneID | attrCmnextCloneRefCnt,
	}

	var buf [32]byte
	_, _, errno := syscall.Syscall6(
		sysGetattrlist,
		uintptr(unsafe.Pointer(pathPtr)),
		uintptr(unsafe.Pointer(&al)),
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(len(buf)),
		uintptr(fsoptAttrCmnExtended),
		0,
	)
	if errno != 0 {
		return 0, 0, 0, false
	}

	total := binary.LittleEndian.Uint32(buf[0:4])
	if total < 20 {
		return 0, 0, 0, false
	}

	docID = binary.LittleEndian.Uint32(buf[4:8])
	cloneID = binary.LittleEndian.Uint64(buf[8:16])
	refCount = binary.LittleEndian.Uint32(buf[16:20])
	return docID, cloneID, refCount, true
}