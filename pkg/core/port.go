package core

import (
	"bufio"
	"context"
	"encoding/binary"
	"net"

	pb "gl2/pkg/pb"
	"gl2/pkg/metrics"
)

type Port struct {
	conn       net.Conn
	udpConn    *net.UDPConn
	udpAddr    *net.UDPAddr
	grpcStream pb.TunnelService_StreamServer
	send       chan []byte
	ctx        context.Context
	cancel     context.CancelFunc
	mode       string
	metrics    *metrics.Metrics
	pool       *BufferPool
}

func NewPort(conn net.Conn, m *metrics.Metrics, pool *BufferPool) *Port {
	ctx, cancel := context.WithCancel(context.Background())
	p := &Port{
		conn:    conn,
		send:    make(chan []byte, 4096),
		ctx:     ctx,
		cancel:  cancel,
		mode:    "stream",
		metrics: m,
		pool:    pool,
	}
	go p.writeLoopStream()
	return p
}

func NewUDPPort(udpConn *net.UDPConn, addr *net.UDPAddr, m *metrics.Metrics, pool *BufferPool) *Port {
	ctx, cancel := context.WithCancel(context.Background())
	p := &Port{
		udpConn: udpConn,
		udpAddr: addr,
		send:    make(chan []byte, 4096),
		ctx:     ctx,
		cancel:  cancel,
		mode:    "udp",
		metrics: m,
		pool:    pool,
	}
	go p.writeLoopUDP()
	return p
}

func NewGRPCPort(stream pb.TunnelService_StreamServer, m *metrics.Metrics, pool *BufferPool) *Port {
	ctx, cancel := context.WithCancel(stream.Context())
	p := &Port{
		grpcStream: stream,
		send:       make(chan []byte, 4096),
		ctx:        ctx,
		cancel:     cancel,
		mode:       "grpc",
		metrics:    m,
		pool:       pool,
	}
	go p.writeLoopGRPC()
	return p
}

func (p *Port) writeLoopStream() {
	defer p.conn.Close()
	writer := bufio.NewWriterSize(p.conn, 64*1024)

	for {
		select {
		case frame, ok := <-p.send:
			if !ok {
				return
			}
			if _, err := writer.Write(frame); err != nil {
				p.pool.Put(frame)
				return
			}
			p.metrics.IncTxBytes(uint64(len(frame)))
			p.metrics.IncTxPackets(1)
			p.pool.Put(frame)
			if len(p.send) == 0 {
				if err := writer.Flush(); err != nil {
					return
				}
			}
		case <-p.ctx.Done():
			return
		}
	}
}

func (p *Port) writeLoopUDP() {
	for {
		select {
		case frame, ok := <-p.send:
			if !ok {
				return
			}
			_, _ = p.udpConn.WriteToUDP(frame, p.udpAddr)
			p.metrics.IncTxBytes(uint64(len(frame)))
			p.metrics.IncTxPackets(1)
			p.pool.Put(frame)
		case <-p.ctx.Done():
			return
		}
	}
}

func (p *Port) writeLoopGRPC() {
	for {
		select {
		case frame, ok := <-p.send:
			if !ok {
				return
			}
			err := p.grpcStream.Send(&pb.Frame{Payload: frame})
			p.metrics.IncTxBytes(uint64(len(frame)))
			p.metrics.IncTxPackets(1)
			p.pool.Put(frame)
			if err != nil {
				return
			}
		case <-p.ctx.Done():
			return
		}
	}
}

func (p *Port) Queue(frame []byte) {
	buf := p.pool.Get()
	buf = buf[:cap(buf)]
	length := len(frame)

	if p.mode == "udp" || p.mode == "grpc" {
		if length > len(buf) {
			buf = make([]byte, length)
		}
		copy(buf[:length], frame)
		select {
		case p.send <- buf[:length]:
		default:
			p.pool.Put(buf)
		}
	} else {
		if length+4 > len(buf) {
			buf = make([]byte, length+4)
		}
		binary.BigEndian.PutUint32(buf[:4], uint32(length))
		copy(buf[4:4+length], frame)
		select {
		case p.send <- buf[:4+length]:
		default:
			p.pool.Put(buf)
		}
	}
}

func (p *Port) Close() {
	p.cancel()
	close(p.send)
	if p.conn != nil {
		p.conn.Close()
	}
}
