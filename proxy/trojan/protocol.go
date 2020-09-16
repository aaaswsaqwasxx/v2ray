package trojan

import (
	"encoding/binary"
	"io"
	"sync"

	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/protocol"
)

var (
	crlf = []byte{'\r', '\n'}

	addrParser = protocol.NewAddressParser(
		protocol.AddressFamilyByte(0x01, net.AddressFamilyIPv4),
		protocol.AddressFamilyByte(0x04, net.AddressFamilyIPv6),
		protocol.AddressFamilyByte(0x03, net.AddressFamilyDomain),
	)
)

const (
	maxLength = 8192

	CommandTCP byte = 1
	CommandUDP byte = 3
)

type ConnWriter struct {
	io.Writer
	Target     net.Destination
	Account    *MemoryAccount
	headerSent bool
}

func (c *ConnWriter) Write(p []byte) (n int, err error) {
	if !c.headerSent {
		if err := c.writeHeader(); err != nil {
			return 0, newError("failed to write request header").Base(err)
		}
	}

	return c.Writer.Write(p)
}

func (c *ConnWriter) WriteMultiBuffer(mb buf.MultiBuffer) error {
	defer buf.ReleaseMulti(mb)

	for _, b := range mb {
		if !b.IsEmpty() {
			if _, err := c.Write(b.Bytes()); err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *ConnWriter) writeHeader() error {
	buffer := buf.StackNew()
	defer buffer.Release()

	command := CommandTCP
	if c.Target.Network == net.Network_UDP {
		command = CommandUDP
	}

	buffer.Write(c.Account.Key)
	buffer.Write(crlf)
	buffer.WriteByte(command)
	addrParser.WriteAddressPort(&buffer, c.Target.Address, c.Target.Port)
	buffer.Write(crlf)

	_, err := c.Writer.Write(buffer.Bytes())
	if err == nil {
		c.headerSent = true
	}

	return err
}

type PacketWriter struct {
	sync.Mutex
	io.Writer
	Target net.Destination
}

func (w *PacketWriter) WriteMultiBuffer(mb buf.MultiBuffer) error {
	w.Lock()
	defer w.Unlock()

	b := make([]byte, maxLength)
	for !mb.IsEmpty() {
		var length int
		mb, length = buf.SplitBytes(mb, b)
		if _, err := w.writePacket(b[:length]); err != nil {
			buf.ReleaseMulti(mb)
			return err
		}
	}

	return nil
}

func (w *PacketWriter) writePacket(payload []byte) (int, error) {
	buffer := buf.StackNew()
	defer buffer.Release()

	length := len(payload)
	lengthBuf := [2]byte{}
	binary.BigEndian.PutUint16(lengthBuf[:], uint16(length))
	addrParser.WriteAddressPort(&buffer, w.Target.Address, w.Target.Port)
	buffer.Write(lengthBuf[:])
	buffer.Write(crlf)
	buffer.Write(payload)
	_, err := w.Write(buffer.Bytes())
	if err != nil {
		return 0, err
	}

	return length, nil
}

type ConnReader struct {
	io.Reader
	Target       net.Destination
	headerParsed bool
}

func (c *ConnReader) ParseHeader() error {
	var crlf [2]byte
	var command [1]byte
	var hash [56]byte
	if _, err := io.ReadFull(c.Reader, hash[:]); err != nil {
		return newError("failed to read user hash").Base(err)
	}

	if _, err := io.ReadFull(c.Reader, crlf[:]); err != nil {
		return newError("failed to read crlf").Base(err)
	}

	if _, err := io.ReadFull(c.Reader, command[:]); err != nil {
		return newError("failed to read command").Base(err)
	}

	network := net.Network_TCP
	if command[0] == CommandUDP {
		network = net.Network_UDP
	}

	addr, port, err := addrParser.ReadAddressPort(nil, c.Reader)
	if err != nil {
		return newError("failed to read address and port").Base(err)
	}
	c.Target = net.Destination{Network: network, Address: addr, Port: port}

	if _, err := io.ReadFull(c.Reader, crlf[:]); err != nil {
		return newError("failed to read crlf").Base(err)
	}

	c.headerParsed = true
	return nil
}

func (c *ConnReader) Read(p []byte) (int, error) {
	if !c.headerParsed {
		if err := c.ParseHeader(); err != nil {
			return 0, err
		}
	}

	return c.Reader.Read(p)
}

func (c *ConnReader) ReadMultiBuffer() (buf.MultiBuffer, error) {
	b := buf.New()
	_, err := b.ReadFrom(c)
	return buf.MultiBuffer{b}, err
}

type PacketPayload struct {
	Target net.Destination
	Buffer buf.MultiBuffer
}

type PacketReader struct {
	io.Reader
}

func (r *PacketReader) ReadMultiBuffer() (buf.MultiBuffer, error) {
	p, err := r.ReadPacket()
	if p != nil {
		return p.Buffer, err
	}
	return nil, err
}

func (r *PacketReader) ReadPacket() (*PacketPayload, error) {
	addr, port, err := addrParser.ReadAddressPort(nil, r)
	if err != nil {
		return nil, newError("failed to read address and port").Base(err)
	}

	var lengthBuf [2]byte
	if _, err := io.ReadFull(r, lengthBuf[:]); err != nil {
		return nil, newError("failed to read payload length").Base(err)
	}

	remain := int(binary.BigEndian.Uint16(lengthBuf[:]))
	if remain > maxLength {
		return nil, newError("oversize payload")
	}

	var crlf [2]byte
	if _, err := io.ReadFull(r, crlf[:]); err != nil {
		return nil, newError("failed to read crlf").Base(err)
	}

	dest := net.UDPDestination(addr, port)
	var mb buf.MultiBuffer
	for remain > 0 {
		length := buf.Size
		if remain < length {
			length = remain
		}

		b := buf.New()
		mb = append(mb, b)
		n, err := b.ReadFullFrom(r, int32(length))
		if err != nil {
			return &PacketPayload{Target: dest, Buffer: mb}, newError("failed to read payload").Base(err)
		}

		remain -= int(n)
	}

	return &PacketPayload{Target: dest, Buffer: mb}, nil
}

// MultiBufferContainer is a Read wrapper over MultiBuffer.
type MultiBufferContainer struct {
	buf.MultiBuffer
}

// ReadMultiBuffer implements Reader.
func (c *MultiBufferContainer) ReadMultiBuffer() (buf.MultiBuffer, error) {
	if c.MultiBuffer.IsEmpty() {
		return nil, io.EOF
	}

	mb := c.MultiBuffer
	c.MultiBuffer = nil
	return mb, nil
}
