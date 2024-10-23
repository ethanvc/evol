package main

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

type HttpRespPacket struct {
	Status     string
	StatusCode int
	Proto      string
	ProtoMajor int
	ProtoMinor int
	Header     http.Header
	Body       []byte
}
