package entity

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestCache_Set(t *testing.T) {
	type fields struct {
		items     map[string]MemoryCacheItem
		mu        sync.RWMutex
		expiryChs map[string]chan struct{}
		onExpired func(string)
	}
	type args struct {
		key   string
		value interface{}
		ttl   time.Duration
	}

	finish := func(s string) {
		println("expired key:", s)
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "test",
			fields: fields{
				items:     make(map[string]MemoryCacheItem),
				expiryChs: make(map[string]chan struct{}),
				onExpired: finish,
				mu:        sync.RWMutex{},
			},
			args: args{
				key:   "test",
				value: "test",
				ttl:   time.Second * 5,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &MemoryCache{
				items:     tt.fields.items,
				mu:        tt.fields.mu,
				expiryChs: tt.fields.expiryChs,
				onExpired: tt.fields.onExpired,
			}
			c.Set(tt.args.key, tt.args.value, tt.args.ttl)
			// 循环7秒，每秒检查一次
			for i := 0; i < 7; i++ {
				println("i:", i)
				time.Sleep(time.Second)
				v, ok := c.Get(tt.args.key)
				if !ok {
					t.Error("get cache error")
				}
				fmt.Println("v:", v)
			}
		})
	}
}
