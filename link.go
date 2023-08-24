package libbpfgo

/*
#cgo LDFLAGS: -lelf -lz
#include "libbpfgo.h"
*/
import "C"

import (
	"fmt"
	"syscall"
	"unsafe"
)

//
// LinkType
//

type LinkType int

const (
	Tracepoint LinkType = iota
	RawTracepoint
	Kprobe
	Kretprobe
	LSM
	PerfEvent
	Uprobe
	Uretprobe
	Tracing
	XDP
	Cgroup
	CgroupLegacy
	Netns
	Iter
)

//
// BPFLinkLegacy
//

type BPFLinkLegacy struct {
	attachType BPFAttachType
	cgroupDir  string
}

//
// BPFLink
//

type BPFLink struct {
	link      *C.struct_bpf_link
	prog      *BPFProg
	linkType  LinkType
	eventName string
	legacy    *BPFLinkLegacy // if set, this is a fake BPFLink
}

func (l *BPFLink) DestroyLegacy(linkType LinkType) error {
	switch l.linkType {
	case CgroupLegacy:
		return l.prog.DetachCgroupLegacy(
			l.legacy.cgroupDir,
			l.legacy.attachType,
		)
	}
	return fmt.Errorf("unable to destroy legacy link")
}

func (l *BPFLink) Destroy() error {
	if l.legacy != nil {
		return l.DestroyLegacy(l.linkType)
	}
	if ret := C.bpf_link__destroy(l.link); ret < 0 {
		return syscall.Errno(-ret)
	}
	l.link = nil
	return nil
}

func (l *BPFLink) FileDescriptor() int {
	return int(C.bpf_link__fd(l.link))
}

// Deprecated: use BPFLink.FileDescriptor() instead.
func (l *BPFLink) GetFd() int {
	return l.FileDescriptor()
}

func (l *BPFLink) Pin(pinPath string) error {
	path := C.CString(pinPath)
	errC := C.bpf_link__pin(l.link, path)
	C.free(unsafe.Pointer(path))
	if errC != 0 {
		return fmt.Errorf("failed to pin link %s to path %s: %w", l.eventName, pinPath, syscall.Errno(-errC))
	}
	return nil
}

func (l *BPFLink) Unpin(pinPath string) error {
	path := C.CString(pinPath)
	errC := C.bpf_link__unpin(l.link)
	C.free(unsafe.Pointer(path))
	if errC != 0 {
		return fmt.Errorf("failed to unpin link %s from path %s: %w", l.eventName, pinPath, syscall.Errno(-errC))
	}
	return nil
}

//
// BPF Link Reader (low-level)
//

func (l *BPFLink) Reader() (*BPFLinkReader, error) {
	fd, errno := C.bpf_iter_create(C.int(l.FileDescriptor()))
	if fd < 0 {
		return nil, fmt.Errorf("failed to create reader: %w", errno)
	}
	return &BPFLinkReader{
		l:  l,
		fd: int(uintptr(fd)),
	}, nil
}