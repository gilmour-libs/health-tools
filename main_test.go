package main

import (
	"github.com/garyburd/redigo/redis"
	"os"
	"testing"
)

var pool *redis.Pool

func TestBlankHost(t *testing.T) {
	_, err := removeKey(pool, "")
	if err == nil {
		t.Error("Should have raised Exception")
	}
}

func TestHostDelete(t *testing.T) {
	ident := "test-host-ident"

	conn := pool.Get()
	defer conn.Close()

	conn.Do("HSET", GilmourHealthKey, ident, "true")

	n, err := removeKey(pool, ident)
	if err != nil {
		t.Error(err)
	} else if n != 1 {
		t.Error("Should have deleted the host entry")
	}
}

func TestHostDeleteFail(t *testing.T) {
	n, err := removeKey(pool, "missing")
	if err != nil {
		t.Error(err)
	} else if n != 0 {
		t.Error("Should not delete the missing host entry")
	}
}

func TestMain(m *testing.M) {
	pool = newRedisPool("127.0.0.1", 6379)

	status := m.Run()
	os.Exit(status)
}
