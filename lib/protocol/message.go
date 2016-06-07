// Copyright (C) 2014 The Protocol Authors.

//go:generate -command genxdr go run ../../vendor/github.com/calmh/xdr/cmd/genxdr/main.go
//go:generate genxdr -o message_xdr.go message.go

package protocol

import (
	"bytes"
	"crypto/sha256"
	"fmt"
)

var (
	sha256OfEmptyBlock        = sha256.Sum256(make([]byte, BlockSize))
	HelloMessageMagic  uint32 = 0x9F79BC40
)

func (f FileInfo) String() string {
	return fmt.Sprintf("File{Name:%q, Flags:0%o, Modified:%d, Version:%v, Length:%d, Blocks:%v}",
		f.Name, f.Flags, f.Modified, f.Version, f.Length, f.Blocks)
}

func (f FileInfo) IsDeleted() bool {
	return f.Flags&FlagDeleted != 0
}

func (f FileInfo) IsInvalid() bool {
	return f.Flags&FlagInvalid != 0
}

func (f FileInfo) IsDirectory() bool {
	return f.Flags&FlagDirectory != 0
}

func (f FileInfo) IsSymlink() bool {
	return f.Flags&FlagSymlink != 0
}

func (f FileInfo) HasPermissionBits() bool {
	return f.Flags&FlagNoPermBits == 0
}

func (f FileInfo) FileLength() int64 {
	if f.IsDirectory() || f.IsDeleted() {
		return 128
	}
	return f.Length
}

func (f FileInfo) FileName() string {
	return f.Name
}

func (f FileInfoTruncated) String() string {
	return fmt.Sprintf("File{Name:%q, Flags:0%o, Modified:%d, Version:%v, Length:%d}",
		f.Name, f.Flags, f.Modified, f.Version, f.Length)
}

func (f FileInfoTruncated) IsDeleted() bool {
	return f.Flags&FlagDeleted != 0
}

func (f FileInfoTruncated) IsInvalid() bool {
	return f.Flags&FlagInvalid != 0
}

func (f FileInfoTruncated) IsDirectory() bool {
	return f.Flags&FlagDirectory != 0
}

func (f FileInfoTruncated) IsSymlink() bool {
	return f.Flags&FlagSymlink != 0
}

func (f FileInfoTruncated) HasPermissionBits() bool {
	return f.Flags&FlagNoPermBits == 0
}

func (f FileInfoTruncated) FileLength() int64 {
	if f.IsDirectory() || f.IsDeleted() {
		return 128
	}
	return f.Length
}

func (f FileInfoTruncated) FileName() string {
	return f.Name
}

// WinsConflict returns true if "f" is the one to choose when it is in
// conflict with "other".
func (f FileInfo) WinsConflict(other FileInfo) bool {
	// If a modification is in conflict with a delete, we pick the
	// modification.
	if !f.IsDeleted() && other.IsDeleted() {
		return true
	}
	if f.IsDeleted() && !other.IsDeleted() {
		return false
	}

	// The one with the newer modification time wins.
	if f.Modified > other.Modified {
		return true
	}
	if f.Modified < other.Modified {
		return false
	}

	// The modification times were equal. Use the device ID in the version
	// vector as tie breaker.
	return f.Version.Compare(other.Version) == ConcurrentGreater
}

func (b BlockInfo) String() string {
	return fmt.Sprintf("Block{%d/%d/%x}", b.Offset, b.Length, b.Hash)
}

// IsEmpty returns true if the block is a full block of zeroes.
func (b BlockInfo) IsEmpty() bool {
	return b.Length == BlockSize && bytes.Equal(b.Hash, sha256OfEmptyBlock[:])
}
