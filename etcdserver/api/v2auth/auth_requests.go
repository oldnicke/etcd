package v2auth

import (
	"context"
	"encoding/json"
	"path"

	"github.com/oldnicke/etcd/etcdserver"
	"github.com/oldnicke/etcd/etcdserver/api/v2error"
	"github.com/oldnicke/etcd/etcdserver/etcdserverpb"

	"go.uber.org/zap"
)

func (s *store) ensureAuthDirectories() error {
	if s.ensuredOnce {
		return nil
	}
	for _, res := range []string{StorePermsPrefix, StorePermsPrefix + "/users/", StorePermsPrefix + "/roles/"} {
		ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
		pe := false
		rr := etcdserverpb.Request{
			Method:    "PUT",
			Path:      res,
			Dir:       true,
			PrevExist: &pe,
		}
		_, err := s.server.Do(ctx, rr)
		cancel()
		if err != nil {
			if e, ok := err.(*v2error.Error); ok {
				if e.ErrorCode == v2error.EcodeNodeExist {
					continue
				}
			}
			if s.lg != nil {
				s.lg.Warn(
					"failed to create auth directories",
					zap.Error(err),
				)
			} else {
				plog.Errorf("failed to create auth directories in the store (%v)", err)
			}
			return err
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()
	pe := false
	rr := etcdserverpb.Request{
		Method:    "PUT",
		Path:      StorePermsPrefix + "/enabled",
		Val:       "false",
		PrevExist: &pe,
	}
	_, err := s.server.Do(ctx, rr)
	if err != nil {
		if e, ok := err.(*v2error.Error); ok {
			if e.ErrorCode == v2error.EcodeNodeExist {
				s.ensuredOnce = true
				return nil
			}
		}
		return err
	}
	s.ensuredOnce = true
	return nil
}

func (s *store) enableAuth() error {
	_, err := s.updateResource("/enabled", true)
	return err
}
func (s *store) disableAuth() error {
	_, err := s.updateResource("/enabled", false)
	return err
}

func (s *store) detectAuth() bool {
	if s.server == nil {
		return false
	}
	value, err := s.requestResource("/enabled", false)
	if err != nil {
		if e, ok := err.(*v2error.Error); ok {
			if e.ErrorCode == v2error.EcodeKeyNotFound {
				return false
			}
		}
		if s.lg != nil {
			s.lg.Warn(
				"failed to detect auth settings",
				zap.Error(err),
			)
		} else {
			plog.Errorf("failed to detect auth settings (%s)", err)
		}
		return false
	}

	var u bool
	err = json.Unmarshal([]byte(*value.Event.Node.Value), &u)
	if err != nil {
		if s.lg != nil {
			s.lg.Warn(
				"internal bookkeeping value for enabled isn't valid JSON",
				zap.Error(err),
			)
		} else {
			plog.Errorf("internal bookkeeping value for enabled isn't valid JSON (%v)", err)
		}
		return false
	}
	return u
}

func (s *store) requestResource(res string, quorum bool) (etcdserver.Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()
	p := path.Join(StorePermsPrefix, res)
	method := "GET"
	if quorum {
		method = "QGET"
	}
	rr := etcdserverpb.Request{
		Method: method,
		Path:   p,
		Dir:    false, // TODO: always false?
	}
	return s.server.Do(ctx, rr)
}

func (s *store) updateResource(res string, value interface{}) (etcdserver.Response, error) {
	return s.setResource(res, value, true)
}
func (s *store) createResource(res string, value interface{}) (etcdserver.Response, error) {
	return s.setResource(res, value, false)
}
func (s *store) setResource(res string, value interface{}, prevexist bool) (etcdserver.Response, error) {
	err := s.ensureAuthDirectories()
	if err != nil {
		return etcdserver.Response{}, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()
	data, err := json.Marshal(value)
	if err != nil {
		return etcdserver.Response{}, err
	}
	p := path.Join(StorePermsPrefix, res)
	rr := etcdserverpb.Request{
		Method:    "PUT",
		Path:      p,
		Val:       string(data),
		PrevExist: &prevexist,
	}
	return s.server.Do(ctx, rr)
}

func (s *store) deleteResource(res string) error {
	err := s.ensureAuthDirectories()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()
	pex := true
	p := path.Join(StorePermsPrefix, res)
	_, err = s.server.Do(ctx, etcdserverpb.Request{
		Method:    "DELETE",
		Path:      p,
		PrevExist: &pex,
	})
	return err
}
