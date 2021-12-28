package zookeeper

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-zookeeper/zk"
	"path/filepath"
	"strings"
)

type Option func(o *options)

type options struct {
	ctx  context.Context
	path string
}

//  WithContext with registry context.
func WithContext(ctx context.Context) Option {
	return func(o *options) {
		o.ctx = ctx
	}
}

// WithPath is config path
func WithPath(p string) Option {
	return func(o *options) {
		o.path = p
	}
}

type source struct {
	options *options
	client  *zk.Conn
}

func New(conn *zk.Conn, opts ...Option) (config.Source, error) {
	options := &options{
		ctx:  context.Background(),
		path: "",
	}

	for _, opt := range opts {
		opt(options)
	}

	if options.path == "" {
		return nil, errors.New("path invalid")
	}

	return &source{
		client:  conn,
		options: options,
	}, nil
}

// Load return the config values
func (s *source) Load() ([]*config.KeyValue, error) {
	serviceNamePath := s.options.path
	data, _, err := s.client.Get(serviceNamePath)
	if err != nil {
		fmt.Printf("查询%s失败, err: %v\n", serviceNamePath, err)
		return nil, err
	}

	kvs := make([]*config.KeyValue, 0)
	kvs = append(kvs, &config.KeyValue{
		Key:    serviceNamePath,
		Value:  data,
		Format: strings.TrimPrefix(filepath.Ext(serviceNamePath), "."),
	})

	return kvs, nil
}

// Watch return the watcher
func (s *source) Watch() (config.Watcher, error) {
	return newWatcher(s), nil
}
