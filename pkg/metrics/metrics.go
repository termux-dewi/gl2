package metrics

import (
	"encoding/json"
	"sync/atomic"
	"time"
)

type Metrics struct {
	RxBytes      atomic.Uint64
	TxBytes      atomic.Uint64
	RxPackets    atomic.Uint64
	TxPackets    atomic.Uint64
	ActivePorts  atomic.Uint64
	Uptime       time.Time
}

type MetricsSnapshot struct {
	RxBytes     uint64        `json:"rx_bytes"`
	TxBytes     uint64        `json:"tx_bytes"`
	RxPackets   uint64        `json:"rx_packets"`
	TxPackets   uint64        `json:"tx_packets"`
	ActivePorts uint64        `json:"active_ports"`
	Uptime      time.Duration `json:"uptime"`
}

func NewMetrics() *Metrics {
	return &Metrics{
		Uptime: time.Now(),
	}
}

func (m *Metrics) IncRxBytes(n uint64) {
	m.RxBytes.Add(n)
}

func (m *Metrics) IncTxBytes(n uint64) {
	m.TxBytes.Add(n)
}

func (m *Metrics) IncRxPackets(n uint64) {
	m.RxPackets.Add(n)
}

func (m *Metrics) IncTxPackets(n uint64) {
	m.TxPackets.Add(n)
}

func (m *Metrics) SetActivePorts(n uint64) {
	m.ActivePorts.Store(n)
}

func (m *Metrics) Snapshot() MetricsSnapshot {
	return MetricsSnapshot{
		RxBytes:     m.RxBytes.Load(),
		TxBytes:     m.TxBytes.Load(),
		RxPackets:   m.RxPackets.Load(),
		TxPackets:   m.TxPackets.Load(),
		ActivePorts: m.ActivePorts.Load(),
		Uptime:      time.Since(m.Uptime),
	}
}

func (s MetricsSnapshot) JSON() []byte {
	data, _ := json.Marshal(s)
	return data
}
