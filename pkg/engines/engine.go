package engines

import (
	"context"
	"fmt"
	"sync"

	"github.com/tai-kun/surrealdb.go/pkg/codec"
)

type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *RPCError) Error() string {
	return fmt.Sprintf("%s (code=%d)", e.Message, e.Code)
}

type NullString struct {
	String string
	Valid  bool
}

type ConnectionInfo struct {
	Endpoint string

	// session
	mu *sync.RWMutex
	ns NullString
	db NullString
	tk NullString
}

type ConnectionInfoSnapshot struct {
	Endpoint  string
	Namespace NullString
	Database  NullString
	Token     NullString
}

func (ci *ConnectionInfo) Namespace() (string, bool) {
	ci.mu.RLock()
	defer ci.mu.RUnlock()
	return ci.ns.String, ci.ns.Valid
}

func (ci *ConnectionInfo) setNS(ns string) {
	ci.ns.String = ns
	ci.ns.Valid = true
}

func (ci *ConnectionInfo) unsetNS() {
	ci.ns.String = ""
	ci.ns.Valid = false
}

func (ci *ConnectionInfo) Database() (string, bool) {
	return ci.db.String, ci.db.Valid
}

func (ci *ConnectionInfo) setDB(db string) {
	ci.db.String = db
	ci.db.Valid = true
}

func (ci *ConnectionInfo) unsetDB() {
	ci.db.String = ""
	ci.db.Valid = false
}

func (ci *ConnectionInfo) Token() (string, bool) {
	ci.mu.RLock()
	defer ci.mu.RUnlock()
	return ci.tk.String, ci.tk.Valid
}

func (ci *ConnectionInfo) setTK(tk string) {
	ci.tk.String = tk
	ci.tk.Valid = true
}

func (ci *ConnectionInfo) unsetTK() {
	ci.tk.String = ""
	ci.tk.Valid = false
}

func (ci *ConnectionInfo) Snapshot() ConnectionInfoSnapshot {
	ci.mu.RLock()
	defer ci.mu.RUnlock()
	return ConnectionInfoSnapshot{
		Endpoint:  ci.Endpoint,
		Namespace: ci.ns,
		Database:  ci.db,
		Token:     ci.tk,
	}
}

type Engine interface {
	ConnectionInfo() ConnectionInfo
	Connect(ctx context.Context, endpoint string) error
	Close(ctx context.Context) error
	Send(ctx context.Context, dst any, method string, params []any) error
}

func clone(f codec.Formatter, v any) (any, error) {
	data, err := f.Marshal(v)
	if err != nil {
		err := fmt.Errorf("failed to clone: failed to marshal value: %w", err)
		return nil, err
	}

	if err := f.Unmarshal(data, v); err != nil {
		err := fmt.Errorf("failed to clone: failed to unmarshal data: %w", err)
		return nil, err
	}

	return v, nil
}
