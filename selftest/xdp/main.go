package main

import "C"

import (
	"errors"
	"fmt"
	"syscall"

	bpf "github.com/aquasecurity/libbpfgo"
)

func main() {

	// supported, innerErr := bpf.BPFHelperIsSupported(bpf.BPFProgTypeKprobe, bpf.BPFFuncKtimeGetBootNs)
	// supported, innerErr := bpf.BPFHelperIsSupported(bpf.BPFProgTypeLsm, bpf.BPFFuncKtimeGetBootNs)
	supported, innerErr := bpf.BPFHelperIsSupported(bpf.BPFProgTypeLsm, bpf.BPFFuncGetCurrentCgroupId)
	// supported, innerErr := bpf.BPFHelperIsSupported(bpf.BPFProgTypeKprobe, bpf.BPFFuncKtimeGetNs)

	fmt.Printf("supported: %v, innerErr: %v\n", supported, innerErr)

	// only report if operation not permitted
	if errors.Is(innerErr, syscall.EPERM) {
		fmt.Printf("only report if operation not permitted - supported: %v, innerErr: %v\n", supported, innerErr)
	}

}
