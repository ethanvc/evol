package mitm

import (
	"context"
	"log/slog"
	"net/http"
	"net/url"
)

type PacketStore struct {
	ch chan *HttpPacket
}

func NewPacketStore() *PacketStore {
	s := &PacketStore{
		ch: make(chan *HttpPacket, 1000),
	}
	go s.consumePacket()
	return s
}

func (store *PacketStore) AppendHttpPacket(packet *HttpPacket) {
	store.ch <- packet
}

func (store *PacketStore) consumePacket() {
	for packet := range store.ch {
		slog.InfoContext(context.Background(), "REQ_END", slog.String("url", packet.Req.Url.String()))
	}
}

type HttpPacket struct {
	Req  HttpReqPacket
	Resp HttpRespPacket
	Err  error
}

type HttpReqPacket struct {
	Method     string
	Url        *url.URL
	Proto      string
	ProtoMajor int
	ProtoMinor int
	Header     http.Header
	Body       []byte
}

func (packet *HttpReqPacket) Build(req *http.Request) {
	packet.Method = req.Method
	packet.Url = req.URL
	packet.Proto = req.Proto
	packet.ProtoMajor = req.ProtoMajor
	packet.ProtoMinor = req.ProtoMinor
	packet.Header = req.Header
}

type HttpRespPacket struct {
	Status     string
	StatusCode int
	Proto      string
	ProtoMajor int
	ProtoMinor int
	Header     http.Header
	Body       []byte
}

func (packet *HttpRespPacket) Build(resp *http.Response) {
	packet.Status = resp.Status
	packet.StatusCode = resp.StatusCode
	packet.Proto = resp.Proto
	packet.ProtoMajor = resp.ProtoMajor
	packet.ProtoMinor = resp.ProtoMinor
	packet.Header = resp.Header
}
