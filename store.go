package opplasmairys

import (
	"context"
	"io"
	"log/slog"
)

var _ KVStore = (*DAStore)(nil)

type KVStore interface {
	// Get retrieves the given key if it's present in the key-value data store.
	Get(ctx context.Context, key []byte) ([]byte, error)
	// Put inserts the given value into the key-value data store.
	Put(ctx context.Context, key []byte, value []byte) error
}

type DAStore struct {
	client *IrysClient
	cache  KVStore
}

func NewDAStore(c Config, store KVStore) (*DAStore, error) {
	client, err := NewIrysClient(c.NetworkName, c.NetWorkRpc, c.PrivateKey, c.FreeUpload)
	if err != nil {
		return nil, err
	}
	return &DAStore{
		client: client,
		cache:  store,
	}, nil
}

func (d *DAStore) Get(ctx context.Context, key []byte) (v []byte, err error) {
	if d.cache != nil {
		v, err = d.cache.Get(ctx, key)
		if err != nil {
			return nil, err
		}
	}
	if v != nil {
		return v, nil
	}
	r, err := d.client.Download(ctx, key)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	v, err = io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	if d.cache != nil {
		err1 := d.cache.Put(ctx, key, v)
		if err1 != nil {
			slog.Warn("failed to cache value, ", "err", err1)
		}
	}
	return v, nil

}

func (d *DAStore) Put(ctx context.Context, key []byte, value []byte) error {
	if d.cache != nil {
		if err := d.cache.Put(ctx, key, value); err != nil {
			return err
		}
	}
	return d.client.Upload(ctx, key, value)
}
