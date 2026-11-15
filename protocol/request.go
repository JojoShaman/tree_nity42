package protocol

type Base struct {
	PID uint32
}

type InfoReq struct {
	Client string
	Base
}

type DisconnectReq struct {
	Client string
	Base
}

type HeartbeatReq struct {
	Client string
	Base
}

type ListReq struct {
	Base
}

type CreateReq struct {
	Topic string
	Base
}

type ProduceReq struct {
	Topic string
	Key   string
	Value string
	Base
}
type CommitReq struct {
	Client string
	Offset uint32
	Base
}

type SubscribeReq struct {
	Topic      string
	Name       string
	Prefix     string
	Offset     uint32
	HassOffset bool
	Base
}

type Request interface {
	Type() byte
	GetPID() uint32
}

func (b Base) GetPID() uint32 { return b.PID }

func (r InfoReq) Type() byte       { return TypeInfo }
func (r ListReq) Type() byte       { return TypeList }
func (r CreateReq) Type() byte     { return TypeCreate }
func (r CommitReq) Type() byte     { return TypeCreate }
func (r ProduceReq) Type() byte    { return TypeProduce }
func (r SubscribeReq) Type() byte  { return TypeSubscribe }
func (r HeartbeatReq) Type() byte  { return TypeHeartbeat }
func (r DisconnectReq) Type() byte { return TypeDisconnect }

const (
	TypeInfo byte = iota + 1
	TypeList
	TypeCommit
	TypeCreate
	TypeProduce
	TypeSubscribe
	TypeDisconnect
	TypeHeartbeat
)
