package gogoing

import (
	"io"
	"bufio"
	"bytes"
	"sync"
	"encoding/binary"
	"errors"
	"fmt"
)

const HeaderSize  = 2

var (
	dataPacketTooBig = errors.New("ReadStream: data packet too big")
	invalidDataPacket = errors.New("ReadStream: invalid packet data")
	dataPacketEncodeError = errors.New("WrireStream: encode stream error")
)


type EventStream interface {
	Read()(*Event, error)
	Write(*Event) error
	Flush() error
	Close() error
	Raw() io.ReadWriteCloser
	MaxPacketSize(size int)
}

type evStream struct {

	streamMutex sync.RWMutex
	conn io.ReadWriteCloser

	outputWriter *bufio.Writer
	outputHeadBuffer *bytes.Buffer

	inputHeadBuffer []byte
	inputHeadReader *bytes.Reader

	encoder Encoder
	decoder Decoder

	maxPacketSize int
}

func (self *evStream) MaxPacketSize(size int) {
	self.maxPacketSize = size
}

func (self *evStream) Read() (e *Event, err error) {
	if _, err = self.inputHeadReader.Seek(0, 0); err != nil { // headReader重置读取位置
		return nil, err
	}
	if _, err = io.ReadFull(self.conn, self.inputHeadBuffer); err != nil { // 读取header
		return nil, err
	}

	// 读取包大小
	var fullsize uint16
	if err = binary.Read(self.inputHeadReader, binary.LittleEndian, &fullsize); err != nil {
		return nil, err
	}

	if self.maxPacketSize>0 && int(fullsize)>self.maxPacketSize {
		return nil, dataPacketTooBig
	}
	eventDataSize := fullsize - HeaderSize
	if eventDataSize<0 {
		return nil, invalidDataPacket
	}
	data := make([]byte, eventDataSize)
	if _, err = io.ReadFull(self.conn, data); err != nil {
		return nil, err
	}
	e, err = self.decoder.Decode(data)
	return
}

func (self *evStream) Write(e *Event) (err error) {
	self.streamMutex.Lock()
	defer self.streamMutex.Unlock()

	self.outputHeadBuffer.Reset()

	// 包体编码
	data, encodeErr := self.encoder.Encode(e)
	if encodeErr!=nil {
		 return encodeErr
	}
	if data == nil {
		return dataPacketEncodeError
	}

	// 包头数据
	if err = binary.Write(self.outputHeadBuffer, binary.LittleEndian, uint16(len(data)+HeaderSize)); err!=nil {
		return err
	}

	// 发送包头
	if err = self.writeFull(self.outputHeadBuffer.Bytes()); err !=nil {
		return err
	}

	// 发送包体
	if err = self.writeFull(data); err !=nil {
		return err
	}
	return

}

const sendTotalTryCount  = 100

func (self *evStream) Flush()(err error) {
	for tryTimes:=0;tryTimes<100;tryTimes++  {
		err = self.outputWriter.Flush()

		// 如果没写完, flush底层会将没发完的buff准备好,
		if err != io.ErrShortWrite {
			break
		}
	}
	return err
}

func (self *evStream) Close() (err error) {
	return self.conn.Close()
}

func (self *evStream) Raw() io.ReadWriteCloser {
	return self.conn
}

func (self *evStream) writeFull(data []byte) error {
	size := len(data)
	for {
		n, err := self.outputWriter.Write(data)
		if err != nil {
			return err
		}
		if n>=size {
			break
		}
		// 数据大于发送缓存
		copy(data[0: size-n], data[n:size])
		size -= n
	}
	return nil
}

func NewStreamWithCoder(conn io.ReadWriteCloser, encoder Encoder, decoder Decoder) (stream *evStream) {

	stream = &evStream{
		conn:conn,
		outputWriter:bufio.NewWriter(conn),
		outputHeadBuffer:bytes.NewBuffer([]byte{}),
		inputHeadBuffer:make([]byte, HeaderSize),
		encoder:encoder,
		decoder:decoder,
	}
	stream.inputHeadReader = bytes.NewReader(stream.inputHeadBuffer)
	return
}


func NewStream(conn io.ReadWriteCloser) (stream *evStream) {

	stream = &evStream{
		conn:conn,
		outputWriter:bufio.NewWriter(conn),
		outputHeadBuffer:bytes.NewBuffer([]byte{}),
		inputHeadBuffer:make([]byte, HeaderSize),
		encoder:NewDefaultEncoder(),
		decoder:NewDefaultDecoder(),
	}
	stream.inputHeadReader = bytes.NewReader(stream.inputHeadBuffer)
	return
}
