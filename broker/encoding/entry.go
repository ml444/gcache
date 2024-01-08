package encoding

import (
	"encoding/binary"
	"errors"
)

const (
	ItemHashByteLen      = 8
	ItemKeySizeByteLen   = 4
	ItemValueSizeByteLen = 4
	ItemIndexByteLen     = ItemHashByteLen + ItemKeySizeByteLen + ItemValueSizeByteLen
)

type Entry struct {
	Key        string
	Value      []byte
	keyLen     uint32
	valueLen   uint32
	DataOffset uint64
}

func (e *Entry) GetKeyLen() uint32 {
	return e.keyLen
}
func (e *Entry) GetValueLen() uint32 {
	return e.valueLen
}
func (e *Entry) GetDataEnd() uint64 {
	return e.DataOffset + uint64(e.keyLen+e.valueLen)
}

type Encode struct {
	Endian binary.ByteOrder
}

func (e *Encode) Marshal(entry *Entry) ([]byte, []byte, error) {
	keyLen := len(entry.Key)
	valueLen := len(entry.Value)
	indexBuf := make([]byte, ItemIndexByteLen)
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
	indexBuf := make([]byte, ItemIndexByteLen)
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
	if len(buf) < ItemIndexByteLen {
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
