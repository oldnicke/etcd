// Package httptypes defines how etcd's HTTP API entities are serialized to and
// deserialized from JSON.
package httptypes

import (
	"encoding/json"

	"go.etcd.io/etcd/pkg/types"
)

type Member struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	PeerURLs   []string `json:"peerURLs"`
	ClientURLs []string `json:"clientURLs"`
}

type MemberCreateRequest struct {
	PeerURLs types.URLs
}

type MemberUpdateRequest struct {
	MemberCreateRequest
}

func (m *MemberCreateRequest) UnmarshalJSON(data []byte) error {
	s := struct {
		PeerURLs []string `json:"peerURLs"`
	}{}

	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	urls, err := types.NewURLs(s.PeerURLs)
	if err != nil {
		return err
	}

	m.PeerURLs = urls
	return nil
}

type MemberCollection []Member

func (c *MemberCollection) MarshalJSON() ([]byte, error) {
	d := struct {
		Members []Member `json:"members"`
	}{
		Members: []Member(*c),
	}

	return json.Marshal(d)
}
