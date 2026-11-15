package protocol

import (
	"encoding/binary"
	"fmt"
)

func decodeProduce(body []byte, pid uint32) (ProduceReq, error) {
	r := Reader{body: body}
	topic := r.readString()
	key := r.readString()
	value := r.readString()
	if r.err != nil {
		return ProduceReq{}, r.err
	}
	return ProduceReq{
		Topic: topic,
		Key:   key,
		Value: value,
		Base:  Base{PID: pid},
	}, nil
}

func decodeCommit(body []byte, pid uint32) (CommitReq, error) {
	r := Reader{body: body}
	client := r.readString()
	offset := r.readUint32()
	if r.err != nil {
		return CommitReq{}, r.err
	}
	return CommitReq{
		Client: client,
		Offset: offset,
		Base:   Base{PID: pid},
	}, nil
}

func decodeSubscribe(body []byte, pid uint32) (SubscribeReq, error) {
	r := Reader{body: body}
	flag := r.readByte() == 1
	offset := r.readUint32()
	topic := r.readString()
	name := r.readString()
	prefix := r.readString()
	if r.err != nil {
		return SubscribeReq{}, r.err
	}
	return SubscribeReq{
		Topic:      topic,
		Name:       name,
		Prefix:     prefix,
		Offset:     offset,
		HassOffset: flag,
		Base:       Base{PID: pid},
	}, nil
}

func Decode(header, body []byte) (Request, error) {
	h := header[0]
	pid := binary.LittleEndian.Uint32(header[1:5])
	switch h {
	case TypeCommit:
		return decodeCommit(body, pid)
	case TypeProduce:
		return decodeProduce(body, pid)
	case TypeSubscribe:
		return decodeSubscribe(body, pid)
	case TypeInfo:
		return InfoReq{
			Client: string(body),
			Base:   Base{PID: pid}}, nil
	case TypeCreate:
		return CreateReq{
			Topic: string(body),
			Base:  Base{PID: pid}}, nil
	case TypeList:
		return ListReq{
			Base: Base{PID: pid}}, nil
	case TypeHeartbeat:
		return HeartbeatReq{
			Client: string(body),
			Base:   Base{PID: pid}}, nil
	case TypeDisconnect:
		return DisconnectReq{
			Client: string(body),
			Base:   Base{PID: pid}}, nil
	}
	return nil, fmt.Errorf("unknown request type: %d", h)
}
