package util

import (
	"encoding/json"
	ccom "github.com/ceph/go-ceph/common/commands"
	"github.com/ceph/go-ceph/rados"
)

var (
	dfCmd = []byte(`{"prefix":"df","format":"json"}`)
)

type RadosCommander = ccom.RadosCommander

type CephAdmin struct {
	conn RadosCommander
}

func NewFromConn(conn RadosCommander) *CephAdmin {
	return &CephAdmin{conn}
}

func (rsa *CephAdmin) PoolStatus() (*PoolStatusResult, error) {
	res := rsa.rawMonCommand(dfCmd)
	return parsePoolStatus(res)
}

type PoolStatusResult struct {
	Pools []*PoolStatusItem `json:"pools"`
}

type PoolStatusItem struct {
	Name  string `json:"name"`
	Stats Stats  `json:"stats"`
}

type Stats struct {
	PercentUsed float64 `json:"percent_used"`
}

func parsePoolStatus(res Response) (*PoolStatusResult, error) {
	r := &PoolStatusResult{}
	resp := res.NoStatus().Unmarshal(r)
	if resp.Ok() {
		return r, nil
	} else {
		return nil, resp.err
	}
}

func (rsa *CephAdmin) validate() error {
	if rsa.conn == nil {
		return rados.ErrNotConnected
	}
	return nil
}

// rawMgrCommand takes a byte buffer and sends it to the MGR as a command.
// The buffer is expected to contain preformatted JSON.
func (rsa *CephAdmin) rawMgrCommand(buf []byte) Response {
	return RawMgrCommand(rsa.conn, buf)
}

// marshalMgrCommand takes an generic interface{} value, converts it to JSON and
// sends the json to the MGR as a command.
func (rsa *CephAdmin) marshalMgrCommand(v interface{}) Response {
	return MarshalMgrCommand(rsa.conn, v)
}

// rawMonCommand takes a byte buffer and sends it to the MON as a command.
// The buffer is expected to contain preformatted JSON.
func (rsa *CephAdmin) rawMonCommand(buf []byte) Response {
	return RawMonCommand(rsa.conn, buf)
}

// marshalMonCommand takes an generic interface{} value, converts it to JSON and
// sends the json to the MGR as a command.
func (rsa *CephAdmin) marshalMonCommand(v interface{}) Response {
	return MarshalMonCommand(rsa.conn, v)
}

func validate(m interface{}) error {
	if m == nil {
		return rados.ErrNotConnected
	}
	return nil
}

// RawMgrCommand takes a byte buffer and sends it to the MGR as a command.
// The buffer is expected to contain preformatted JSON.
func RawMgrCommand(m ccom.MgrCommander, buf []byte) Response {
	if err := validate(m); err != nil {
		return Response{err: err}
	}
	return NewResponse(m.MgrCommand([][]byte{buf}))
}

// MarshalMgrCommand takes an generic interface{} value, converts it to JSON
// and sends the json to the MGR as a command.
func MarshalMgrCommand(m ccom.MgrCommander, v interface{}) Response {
	b, err := json.Marshal(v)
	if err != nil {
		return Response{err: err}
	}
	return RawMgrCommand(m, b)
}

// RawMonCommand takes a byte buffer and sends it to the MON as a command.
// The buffer is expected to contain preformatted JSON.
func RawMonCommand(m ccom.MonCommander, buf []byte) Response {
	if err := validate(m); err != nil {
		return Response{err: err}
	}
	return NewResponse(m.MonCommand(buf))
}

// MarshalMonCommand takes an generic interface{} value, converts it to JSON
// and sends the json to the MGR as a command.
func MarshalMonCommand(m ccom.MonCommander, v interface{}) Response {
	b, err := json.Marshal(v)
	if err != nil {
		return Response{err: err}
	}
	return RawMonCommand(m, b)
}
