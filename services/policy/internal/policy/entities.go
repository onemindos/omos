package policy

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/nats-io/nats.go"
)

// Entity is the subset of the entity record (onemind-xr docs/contracts/entity.md) the
// policy service reads.
type Entity struct {
	EntityID   string `json:"entity_id"`
	EntityType string `json:"entity_type"`
	Kind       string `json:"kind"`
	Location   *struct {
		Lat  float64 `json:"lat"`
		Lon  float64 `json:"lon"`
		HaeM float64 `json:"hae_m"`
	} `json:"location"`
	Health *struct {
		Connection string `json:"connection"`
	} `json:"health"`
	Capabilities []string `json:"capabilities"`
	Indicators   *struct {
		Simulated *bool `json:"simulated"`
	} `json:"indicators"`
}

// ErrEntityNotFound means the entity service answered and has no such entity.
var ErrEntityNotFound = errors.New("entity not found")

// EntityResolver looks up one entity. Any error other than ErrEntityNotFound means the
// lookup could not be made, and the policy service denies (fail closed).
type EntityResolver interface {
	Get(ctx context.Context, entityType, id string) (*Entity, error)
}

// NATSEntities calls svc.entity.get (BUS-DESIGN.md section 3.4).
type NATSEntities struct {
	Conn    *nats.Conn
	Timeout time.Duration
}

func (n *NATSEntities) Get(ctx context.Context, entityType, id string) (*Entity, error) {
	req, _ := json.Marshal(map[string]string{"entity_type": entityType, "entity_id": id})
	t := n.Timeout
	if t == 0 {
		t = 2 * time.Second
	}
	ctx, cancel := context.WithTimeout(ctx, t)
	defer cancel()
	msg, err := n.Conn.RequestWithContext(ctx, "svc.entity.get", req)
	if err != nil {
		return nil, err
	}
	var wrap struct {
		Entity *Entity `json:"entity"`
		Error  *struct {
			Code string `json:"code"`
		} `json:"error"`
	}
	if err := json.Unmarshal(msg.Data, &wrap); err != nil {
		return nil, err
	}
	if wrap.Error != nil {
		if wrap.Error.Code == "not_found" {
			return nil, ErrEntityNotFound
		}
		return nil, errors.New("entity service error: " + wrap.Error.Code)
	}
	if wrap.Entity != nil {
		return wrap.Entity, nil
	}
	var e Entity
	if err := json.Unmarshal(msg.Data, &e); err != nil || e.EntityID == "" {
		return nil, errors.New("entity service: unreadable reply")
	}
	return &e, nil
}
