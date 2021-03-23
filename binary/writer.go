package binary

import (
	"encoding/binary"
	"encoding/hex"
	"sync"
)

type Writer struct {
	b []byte
}

var bufferPool = sync.Pool{
	New: func() interface{} {
		return new(Writer)
	},
}

func NewWriter() *Writer {
	return bufferPool.Get().(*Writer)
}

func PutBuffer(w *Writer) {
	// See https://golang.org/issue/23199
	const maxSize = 1 << 16
	if cap(w.b) < maxSize { // 对于大Buffer直接丢弃
		w.b = w.b[:0]
		bufferPool.Put(w)
	}
}

func NewWriterF(f func(writer *Writer)) []byte {
	w := NewWriter()
	f(w)
	b := append([]byte(nil), w.Bytes()...)
	PutBuffer(w)
	return b
}

func (w *Writer) Write(b []byte) {
	w.b = append(w.b, b...)
}

func (w *Writer) WriteHex(h string) {
	b, _ := hex.DecodeString(h)
	w.Write(b)
}

func (w *Writer) WriteByte(b byte) {
	w.b = append(w.b, b)
}

func (w *Writer) WriteUInt16(v uint16) {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, v)
	w.Write(b)
}

func (w *Writer) WriteUInt32(v uint32) {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, v)
	w.Write(b)
}

func (w *Writer) WriteUInt64(v uint64) {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, v)
	w.Write(b)
}

func (w *Writer) WriteString(v string) {
	payload := []byte(v)
	w.WriteUInt32(uint32(len(payload) + 4))
	w.Write(payload)
}

func (w *Writer) WriteStringShort(v string) {
	w.WriteBytesShort([]byte(v))
}

func (w *Writer) WriteBool(b bool) {
	if b {
		w.WriteByte(0x01)
	} else {
		w.WriteByte(0x00)
	}
}

func (w *Writer) EncryptAndWrite(key []byte, data []byte) {
	tea := NewTeaCipher(key)
	ed := tea.Encrypt(data)
	w.Write(ed)
}

func (w *Writer) WriteIntLvPacket(offset int, f func(writer *Writer)) {
	data := NewWriterF(f)
	w.WriteUInt32(uint32(len(data) + offset))
	w.Write(data)
}

func (w *Writer) WriteUniPacket(commandName string, sessionId, extraData, body []byte) {
	w.WriteIntLvPacket(4, func(w *Writer) {
		w.WriteString(commandName)
		w.WriteUInt32(8)
		w.Write(sessionId)
		if len(extraData) == 0 {
			w.WriteUInt32(0x04)
		} else {
			w.WriteUInt32(uint32(len(extraData) + 4))
			w.Write(extraData)
		}
	})
	w.WriteIntLvPacket(4, func(w *Writer) {
		w.Write(body)
	})
}

func (w *Writer) WriteBytesShort(data []byte) {
	w.WriteUInt16(uint16(len(data)))
	w.Write(data)
}

func (w *Writer) WriteTlvLimitedSize(data []byte, limit int) {
	if len(data) <= limit {
		w.WriteBytesShort(data)
		return
	}
	w.WriteBytesShort(data[:limit])
}

func (w *Writer) Bytes() []byte {
	return w.b
}
