package gcache

import (
	"encoding/binary"
	"errors"
)

const (
	itemHashByteLen      = 8
	itemKeySizeByteLen   = 4
	itemValueSizeByteLen = 4
	itemIndexByteLen     = itemHashByteLen + itemKeySizeByteLen + itemValueSizeByteLen
)

type Entry struct {
	Key        string
	Value      []byte
	keyLen     uint32
	valueLen   uint32
	DataOffset uint64
}

type Encode struct {
	Endian binary.ByteOrder
}

func (e *Encode) Marshal(entry *Entry) ([]byte, []byte, error) {
	keyLen := len(entry.Key)
	valueLen := len(entry.Value)
	indexBuf := make([]byte, itemIndexByteLen)
	dataBuf := make([]byte, keyLen+valueLen)
	b := e.Endian
	// index bytes
	b.PutUint32(indexBuf[0:], uint32(keyLen))
	b.PutUint32(indexBuf[4:], uint32(valueLen))
	b.PutUint64(indexBuf[8:], entry.DataOffset)
	// data bytes
	copy(dataBuf[0:], entry.Key)
	copy(dataBuf[keyLen:], entry.Value)
	return indexBuf, dataBuf, nil
}

func (e *Encode) MarshalIndex(entry *Entry) ([]byte, error) {
	keyLen := len(entry.Key)
	valueLen := len(entry.Value)
	indexBuf := make([]byte, itemIndexByteLen)
	b := e.Endian
	// index bytes
	b.PutUint32(indexBuf[0:], uint32(keyLen))
	b.PutUint32(indexBuf[4:], uint32(valueLen))
	b.PutUint64(indexBuf[8:], entry.DataOffset)
	return indexBuf, nil
}

func (e *Encode) MarshalKV(entry *Entry) ([]byte, error) {
	keyLen := len(entry.Key)
	valueLen := len(entry.Value)
	dataBuf := make([]byte, keyLen+valueLen)
	// data bytes
	copy(dataBuf[0:], entry.Key)
	copy(dataBuf[keyLen:], entry.Value)
	return dataBuf, nil
}

func (e *Encode) UnmarshalIndex(buf []byte, entry *Entry) error {
	if len(buf) < itemIndexByteLen {
		return errors.New("invalid index item length")
	}
	//p.dataCorruption = true
	b := e.Endian
	entry.keyLen = b.Uint32(buf[0:])
	entry.valueLen = b.Uint32(buf[4:])
	entry.DataOffset = b.Uint64(buf[8:])
	return nil
}

func (e *Encode) UnmarshalKV(buf []byte, entry *Entry) error {
	keyEnd := entry.keyLen
	entry.Key = string(buf[0:keyEnd])
	entry.Value = buf[keyEnd : keyEnd+entry.valueLen]
	return nil
}
