package modbuslib

import (
	"net"
)

func intTo2Byte(val int) []byte { /* binary.BigEndian.PutUint16 */
	b := make([]byte, 2)
	b[0] = byte(val >> 8)
	b[1] = byte(val)
	return b
}

func byteTo16int(val []byte) int {
	_ = val[1] // bounds check hint to compiler
	return int(val[1]) | int(val[0])<<8
}

type ModbusClient struct {
	Host    string
	Addr    byte
	Code    byte
	Timeout int
	Data    []byte
	Conn    net.Conn
}

func (mb *ModbusClient) Connect() (net.Conn, error) {
	conn, err := net.Dial("tcp", mb.Host)
	mb.Conn = conn
	return conn, err
}

func (mb *ModbusClient) Close() {
	mb.Conn.Close()
}

func (mb *ModbusClient) envelope() []byte {
	if mb.Addr == byte(0) {
		mb.Addr = 0x01
	}

	head := []byte{0x00, 0x00, 0x00, 0x00, 0x00, byte(len(mb.Data) + 2), mb.Addr, mb.Code}
	body := []byte{}
	body = append(body, head...)
	body = append(body, mb.Data...)
	return body
}

func (mb *ModbusClient) ReadHoldingRegister(regAddr int, regSize int) ([]int, error) {
	mb.Code = 0x03
	mb.Data = []byte{}
	mb.Data = append(mb.Data, intTo2Byte(regAddr)...)
	mb.Data = append(mb.Data, intTo2Byte(regSize)...)
	tmp, err := mb.send()
	outp := []int{}

	for i := 0; i < (regSize * 2); i += 2 {
		outp = append(outp, []int{byteTo16int(tmp[9+i : 9+i+2])}...)
	}
	return outp, err
}

func (mb *ModbusClient) send() ([]byte, error) {
	outp := make([]byte, 0x40)

	_, err := mb.Conn.Write(mb.envelope())
	if err != nil {
		return []byte{}, err
	}

	_, err = mb.Conn.Read(outp)
	if err != nil {
		return []byte{}, err
	}

	return outp, err
}
