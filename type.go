package surrealdb

import (
	"fmt"

	"github.com/tai-kun/surrealdb.go/pkg/codec"
)

type Variables = map[string]any

type Auth struct {
	Namespace string
	Database  string
	Access    string
	Username  string
	Password  string
	Variables Variables
}

func (a *Auth) MarshalCBOR() ([]byte, error) {
	return a.marshal(CBORFormatter)
}

func (a *Auth) UnmarshalCBOR(data []byte) error {
	return a.unmarshal(CBORFormatter, data)
}

func (a *Auth) MarshalJSON() ([]byte, error) {
	return a.marshal(JSONFormatter)
}

func (a *Auth) UnmarshalJSON(data []byte) error {
	return a.unmarshal(JSONFormatter, data)
}

func (a *Auth) marshal(f codec.Formatter) ([]byte, error) {
	vars := map[string]any{}
	a.replace(vars)
	data, err := f.Marshal(vars)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Auth: %w", err)
	}
	return data, nil
}

func (a *Auth) unmarshal(f codec.Formatter, data []byte) error {
	var vars map[string]any
	if err := f.Unmarshal(data, &vars); err != nil {
		return fmt.Errorf("failed to unmarshal Auth: %w", err)
	}
	if err := a.revive(vars); err != nil {
		return fmt.Errorf("failed to unmarshal Auth: %w", err)
	}
	return nil
}

func (a *Auth) replace(vars map[string]any) {
	for k, v := range map[string]string{
		"ns":   a.Namespace,
		"db":   a.Database,
		"ac":   a.Access,
		"user": a.Username,
		"pass": a.Password,
	} {
		if v != "" { // omitempty
			vars[k] = v
		}
	}
	for k, v := range a.Variables {
		switch k {
		case "ns", "db", "ac", "user", "pass":
			// omit
		default:
			vars[k] = v
		}
	}
}

func (a *Auth) revive(vars map[string]any) error {
	*a = Auth{}
	for k, v := range vars {
		switch k {
		case "ns", "db", "ac", "user", "pass":
			s, ok := v.(string)
			if !ok {
				err := fmt.Errorf("cannot cast %s (%T) to string", k, v)
				return err
			}
			switch k {
			case "ns":
				a.Namespace = s
			case "db":
				a.Database = s
			case "ac":
				a.Access = s
			case "user":
				a.Username = s
			case "pass":
				a.Password = s
			}
		default:
			a.Variables[k] = v
		}
	}
	return nil
}

func NewRootUserAuth(user, pass string) *Auth {
	return &Auth{
		Username: user,
		Password: pass,
	}
}

func NewNamespaceUserAuth(user, pass, ns string) *Auth {
	return &Auth{
		Username:  user,
		Password:  pass,
		Namespace: ns,
	}
}

func NewDatabaseUserAuth(user, pass, ns, db string) *Auth {
	return &Auth{
		Username:  user,
		Password:  pass,
		Namespace: ns,
		Database:  db,
	}
}

func NewRecordAccessAuth(user, pass, ac string) *Auth {
	return &Auth{
		Username: user,
		Password: pass,
		Access:   ac,
	}
}
