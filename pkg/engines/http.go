package engines

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/fxamacker/cbor/v2"

	"github.com/tai-kun/surrealdb.go/pkg/codec"
	"github.com/tai-kun/surrealdb.go/pkg/models"
)

type httpRPCRequest struct {
	Method string `json:"method"`
	Params []any  `json:"params"`
}

type httpRPCResponse[T any] struct {
	Result T         `json:"result"`
	Error  *RPCError `json:"error"`
}

type HTTPEngine struct {
	mu   sync.RWMutex
	fmt  codec.Formatter
	info *ConnectionInfo
	conn *http.Client
	vars *sync.Map
}

func NewHTTPEngine(fmt codec.Formatter) *HTTPEngine {
	return &HTTPEngine{
		fmt: fmt,
	}
}

func (e *HTTPEngine) ConnectionInfo() ConnectionInfo {
	return *e.info
}

func (e *HTTPEngine) Connect(ctx context.Context, endpoint string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if err := ctx.Err(); err != nil {
		err := fmt.Errorf(
			"engines: http: failed to connect to endpoint %s: %w",
			strconv.Quote(endpoint), err,
		)
		return err
	}

	e.info = &ConnectionInfo{
		Endpoint: endpoint,
		mu:       &e.mu,
	}
	e.conn = &http.Client{
		Timeout: 5 * time.Second,
	}
	e.vars = &sync.Map{}

	return nil
}

func (e *HTTPEngine) Close(ctx context.Context) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	endpoint := e.info.Endpoint
	e.info = nil
	e.conn = nil
	e.vars = nil

	if err := ctx.Err(); err != nil {
		err := fmt.Errorf(
			"engines: http: failed to close from endpoint %s: %w",
			strconv.Quote(endpoint), err,
		)
		return err
	}

	return nil
}

func (e *HTTPEngine) Send(
	ctx context.Context,
	dst any,
	method string,
	params []any,
) error {
	switch method {
	case "use":
		if len(params) != 2 {
			err := fmt.Errorf(
				"engines: http: use: invalid params: needs 2 params, but got %d",
				len(params),
			)
			return err
		}

		info := e.info.Snapshot()
		namespace, database := info.Namespace, info.Database
		ns, db := params[0], params[1]

		switch ns.(type) {
		case models.None, *models.None:
			// pass
		case nil: // Null
			namespace.String = ""
			namespace.Valid = false
		default:
			s, ok := ns.(string)
			if !ok {
				err := fmt.Errorf(
					"engines: http: use: invalid params: "+
						"the namespace should to be a nullish string but got %T",
					ns,
				)
				return err
			}
			namespace.String = s
			namespace.Valid = true
		}

		switch db.(type) {
		case models.None, *models.None:
			// pass
		case nil: // Null
			database.String = ""
			database.Valid = false
		default:
			s, ok := db.(string)
			if !ok {
				err := fmt.Errorf(
					"engines: http: use: invalid params: "+
						"the database should to be a nullish string but got %T",
					db,
				)
				return err
			}
			database.String = s
			database.Valid = true
		}

		if !namespace.Valid && database.Valid {
			err := fmt.Errorf(
				"engines: http: use: missing namespace: "+
					"the namespace must be specified before the database %s",
				strconv.Quote(database.String),
			)
			return err
		}

		e.mu.Lock()
		if namespace.Valid {
			e.info.setNS(namespace.String)
		} else {
			e.info.unsetNS()
		}
		if database.Valid {
			e.info.setDB(database.String)
		} else {
			e.info.unsetDB()
		}
		e.mu.Unlock()

	case "let":
		if len(params) != 1 && len(params) != 2 {
			err := fmt.Errorf(
				"engines: http: let: invalid params: needs 1 or 2 params, but got %d",
				len(params),
			)
			return err
		}
		if len(params) == 1 {
			params = append(params, models.None{})
		}

		p0, v := params[0], params[1]
		k, ok := p0.(string)
		if !ok {
			err := fmt.Errorf(
				"engines: http: let: invalid params: name should to be a string but got %T",
				p0,
			)
			return err
		}

		switch v.(type) {
		case models.None, *models.None:
			e.vars.Delete(k)
		default:
			cloned, err := clone(e.fmt, v)
			if err != nil {
				err := fmt.Errorf("engines: http: let: failed to clone value: %w", err)
				return err
			}

			e.vars.Store(k, cloned)
		}

	case "unset":
		if len(params) == 0 {
			err := fmt.Errorf(
				"engines: http: unset: invalid params: needs a param, but got %d",
				len(params),
			)
			return err
		}

		k, ok := params[0].(string)
		if !ok {
			err := fmt.Errorf(
				"engines: http: unset: invalid params: name should to be a string but got %T",
				params[0],
			)
			return err
		}

		e.vars.Delete(k)

	default:
		info := e.info.Snapshot()
		if !info.Namespace.Valid && info.Database.Valid {
			err := fmt.Errorf(
				"engines: http: %s: missing namespace: "+
					"the namespace must be specified before the database %s",
				method,
				strconv.Quote(info.Database.String),
			)
			return err
		}

		// switch method {
		// case "query":

		// }

		data, err := e.fmt.Marshal(httpRPCRequest{
			Method: method,
			Params: params,
		})
		if err != nil {
			err := fmt.Errorf("engines: http: %s: failed to marshal RPC request: %w", method, err)
			return err
		}

		body := bytes.NewReader(data)
		req, err := http.NewRequestWithContext(ctx, "POST", info.Endpoint, body)
		if err != nil {
			err := fmt.Errorf("engines: http: %s: failed to create a request: %w", method, err)
			return err
		}

		req.Header.Set("Accept", e.fmt.ContentType())
		req.Header.Set("Content-Type", e.fmt.ContentType())
		if info.Namespace.Valid {
			req.Header.Set("Surreal-NS", info.Namespace.String)
		}
		if info.Database.Valid {
			req.Header.Set("Surreal-DB", info.Database.String)
		}
		if info.Token.Valid {
			req.Header.Set("Authorization", "Bearer "+info.Token.String)
		}

		resp, err := e.conn.Do(req)
		if err != nil {
			r := "content-type=" + e.fmt.ContentType()
			if info.Namespace.Valid {
				r += ",ns=" + info.Namespace.String
			}
			if info.Database.Valid {
				r += ",db=" + info.Database.String
			}
			if info.Token.Valid {
				r += ",tk=***"
			}
			err := fmt.Errorf(
				"engines: http: %s: failed to send a request(%s): %w",
				method, r, err,
			)
			return err
		}
		defer resp.Body.Close()

		data, err = io.ReadAll(resp.Body)
		if err != nil {
			err := fmt.Errorf("engines: http: %s: failed to read the response body: %w", method, err)
			return err
		}
		if resp.StatusCode != 200 {
			return fmt.Errorf(
				"engines: http: %s: %d %s: %s",
				method, resp.StatusCode, resp.Status, string(data),
			)
		}

		switch e.fmt.ContentType() {
		case "application/cbor":
			var resp httpRPCResponse[cbor.RawMessage]
			if err := e.fmt.Unmarshal(data, &resp); err != nil {
				err := fmt.Errorf(
					"engines: http: %s: failed to unmarshal RPC response: %w",
					method, err,
				)
				return err
			}
			if resp.Error != nil {
				err := fmt.Errorf("engines: http: %s: failed to execute RPC: %w", method, resp.Error)
				return err
			}
			if err := e.fmt.Unmarshal(resp.Result, dst); err != nil {
				err := fmt.Errorf(
					"engines: http: %s: failed to unmarshal RPC result: %w",
					method, err,
				)
				return err
			}

		default:
			var resp httpRPCResponse[json.RawMessage]
			if err := e.fmt.Unmarshal(data, &resp); err != nil {
				err := fmt.Errorf(
					"engines: http: %s: failed to unmarshal RPC response: %w",
					method, err,
				)
				return err
			}
			if resp.Error != nil {
				err := fmt.Errorf("engines: http: %s: failed to execute RPC: %w", method, resp.Error)
				return err
			}
			if err := e.fmt.Unmarshal(resp.Result, dst); err != nil {
				err := fmt.Errorf(
					"engines: http: %s: failed to unmarshal RPC result: %w",
					method, err,
				)
				return err
			}
		}

		switch method {
		case "signin", "signup":
			s, ok := dst.(*string)
			if !ok {
				msg := fmt.Sprintf(
					"surrealdb: engines: http: %s: invalid response: "+
						"the token should to be a *string but got %T",
					method,
					dst,
				)
				panic(msg)
			}

			e.mu.Lock()
			e.info.setTK(*s)
			e.mu.Unlock()

		case "authenticate":
			if len(params) == 0 {
				msg := fmt.Sprintf(
					"surrealdb: engines: http: authenticate: invalid params: "+
						"needs 1 param, but got %d",
					len(params),
				)
				panic(msg)
			}
			if s, ok := params[0].(string); ok {
				e.mu.Lock()
				e.info.setTK(s)
				e.mu.Unlock()
			} else if s, ok := params[0].(*string); ok {
				e.mu.Lock()
				e.info.setTK(*s)
				e.mu.Unlock()
			} else {
				msg := fmt.Sprintf(
					"surrealdb: engines: http: authenticate: invalid params: "+
						"the token should to be a string but got %T",
					params[0],
				)
				panic(msg)
			}

		case "invalidate":
			e.mu.Lock()
			e.info.unsetTK()
			e.mu.Unlock()
		}
	}

	return nil
}
