package surrealdb

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/fxamacker/cbor/v2"
	"github.com/tai-kun/surrealdb.go/pkg/codec"
	"github.com/tai-kun/surrealdb.go/pkg/engines"
)

type DB struct {
	mu  sync.RWMutex
	ctx context.Context
	eng Engines
	con engines.Engine
	fmt codec.Formatter
}

func New(fmt codec.Formatter, eng Engines) (*DB, error) {
	return &DB{ctx: context.Background(), eng: eng, fmt: fmt}, nil
}

func NewWithContext(
	ctx context.Context,
	eng Engines,
	fmt codec.Formatter,
) (*DB, error) {
	return &DB{ctx: ctx, eng: eng, fmt: fmt}, nil
}

func (db *DB) Connect(endpoint string) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	u, err := processEndpoint(endpoint, transformEndpointDefault)
	if err != nil {
		err = fmt.Errorf("surrealdb: %w", err)
		return err
	}

	endpoint = u.String()
	if db.con != nil {
		connected := db.con.ConnectionInfo().Endpoint
		if endpoint != connected {
			err = fmt.Errorf(
				"surrealdb: an attempt was made to connect to %s while %s was already connected",
				endpoint, connected,
			)
			return err
		}

		return nil
	}

	eng, ok := db.eng[u.Scheme]
	if !ok {
		err = fmt.Errorf("surrealdb: no %s scheme engine found", u.Scheme)
		return err
	}

	con := eng(db.fmt)
	if err := con.Connect(db.ctx, endpoint); err != nil {
		err = fmt.Errorf("surrealdb: %w", err)
		return err
	}
	db.con = con

	return nil
}

func (db *DB) Close() error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if db.con == nil {
		return nil
	}

	defer func() {
		db.con = nil
	}()

	if err := db.con.Close(db.ctx); err != nil {
		err = fmt.Errorf("surrealdb: %w", err)
		return err
	}

	return nil
}

func (db *DB) send(dst any, method string, params ...any) error {
	db.mu.RLock()
	defer db.mu.RUnlock()

	if db.con == nil {
		err := fmt.Errorf("surrealdb: not connected")
		return err
	}

	if err := db.con.Send(db.ctx, dst, method, params); err != nil {
		err := fmt.Errorf("surrealdb: %w", err)
		return err
	}

	return nil
}

func (d *DB) Use(ns, db any) error {
	var r any
	return d.send(&r, "use", ns, db)
}

type CurrentUser = map[string]any

func Info(db *DB) (CurrentUser, error) {
	var r CurrentUser
	if err := db.send(&r, "info"); err != nil {
		return CurrentUser{}, err
	}

	return r, nil
}

func (db *DB) SignUp(auth *Auth) (string, error) {
	var r string
	if err := db.send(&r, "signup", auth); err != nil {
		return "", err
	}

	return r, nil
}

func (db *DB) SignIn(auth *Auth) (string, error) {
	var r string
	if err := db.send(&r, "signin", auth); err != nil {
		return "", err
	}

	return r, nil
}

func (db *DB) Authenticate(token string) (string, error) {
	var r string
	if err := db.send(&r, "authenticate", token); err != nil {
		return "", err
	}

	return r, nil
}

type cborRawQueryResult struct {
	Status string          `json:"status"`
	Time   string          `json:"time"`
	Result cbor.RawMessage `json:"result"`
}

type jsonRawQueryResult struct {
	Status string          `json:"status"`
	Time   string          `json:"time"`
	Result json.RawMessage `json:"result"`
}

type QueryResult struct {
	fmt  codec.Unmarshaler
	data []byte
}

func (qr *QueryResult) Unmarshal(v any) error {
	if err := qr.fmt.Unmarshal(qr.data, v); err != nil {
		err := fmt.Errorf("surrealdb: failed to unmarshal QueryResult: %w", err)
		return err
	}

	return nil
}

type QueryRawResult struct {
	Status string       `json:"status"`
	Time   string       `json:"time"`
	Result *QueryResult `json:"result"`
}

func (db *DB) QueryRaw(surql string, vars Variables) ([]QueryRawResult, error) {
	switch db.fmt.ContentType() {
	case "application/cbor":
		var r1 []cborRawQueryResult
		if err := db.send(&r1, "query", surql, vars); err != nil {
			return nil, err
		}

		r2 := make([]QueryRawResult, len(r1))
		for i, r := range r1 {
			r2[i] = QueryRawResult{
				Status: r.Status,
				Time:   r.Time,
				Result: &QueryResult{
					fmt:  db.fmt,
					data: r.Result,
				},
			}
		}

		return r2, nil

	default:
		var r1 []jsonRawQueryResult
		if err := db.send(&r1, "query", surql, vars); err != nil {
			return nil, err
		}

		r2 := make([]QueryRawResult, len(r1))
		for i, r := range r1 {
			r2[i] = QueryRawResult{
				Status: r.Status,
				Time:   r.Time,
				Result: &QueryResult{
					fmt:  db.fmt,
					data: r.Result,
				},
			}
		}

		return r2, nil
	}
}

type QueryResults struct {
	data []*QueryResult
}

func (qr *QueryResults) Len() int {
	return len(qr.data)
}

func (qr *QueryResults) Remove(i int, v any) error {
	last := qr.Len() - 1
	if i < 0 || i > last {
		return fmt.Errorf(
			"surrealdb: failed to remove QueryResult from QueryResults: index %d out of range",
			i,
		)
	}

	switch i {
	case 0:
		data := qr.data[0]
		qr.data = qr.data[1:]
		return data.Unmarshal(v)

	case last:
		data := qr.data[last]
		qr.data = qr.data[:last]
		return data.Unmarshal(v)

	default:
		data := qr.data[i]
		qr.data = append(qr.data[:i], qr.data[i+1:]...)
		return data.Unmarshal(v)
	}
}

func (db *DB) Query(surql string, vars Variables) (*QueryResults, error) {
	r, err := db.QueryRaw(surql, vars)
	if err != nil {
		return nil, err
	}

	data := make([]*QueryResult, len(r))
	for i, v := range r {
		switch v.Status {
		case "OK":
			data[i] = v.Result

		case "ERR":
			var msg string
			if err := v.Result.Unmarshal(&msg); err != nil {
				panic(err)
			}

			err := fmt.Errorf(
				"surrealdb: failed to execute query at %d of %d statement(s): %s",
				i+1, len(r), msg,
			)
			return nil, err

		default:
			err := fmt.Errorf("surrealdb: unexpected query status %s", v.Status)
			panic(err)
		}
	}

	return &QueryResults{data}, nil
}

func (db *DB) Let(name string, value any) error {
	var r any
	return db.send(&r, "let", name, value)
}

func (db *DB) Unset(name string) error {
	var r any
	return db.send(&r, "unset", name)
}

func (db *DB) Version() (string, error) {
	var r string
	if err := db.send(&r, "version"); err != nil {
		return "", err
	}

	return r, nil
}
