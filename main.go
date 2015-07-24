package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"gopkg.in/gilmour-libs/gilmour-e-go.v0/backends"
	"os"
)

const BLANK = ""
const GilmourHealthKey = "gilmour.known_host.health"

func newRedisPool(host string, port int) *redis.Pool {
	return backends.GetPool(fmt.Sprintf("%v:%v", host, port))
}

func removeKey(pool *redis.Pool, ident string) (int, error) {
	if ident == BLANK {
		err := errors.New("Identifier to remove cannot be blank")
		return 0, err
	}

	conn := pool.Get()
	defer conn.Close()

	return redis.Int(conn.Do("HDEL", GilmourHealthKey, ident))
}

func makeIdent(ident *string) {
	flag.StringVar(ident, "ident", "", "key to cleanup")
}

func makeRedisConf(host *string, port *int) {
	flag.StringVar(host, "host", "127.0.0.1", "Redis host to connect to")
	flag.IntVar(port, "port", 6379, "Redis port")
}

func main() {
	var ident string
	var redis_host string
	var redis_port int

	makeIdent(&ident)
	makeRedisConf(&redis_host, &redis_port)

	flag.Parse()

	pool := newRedisPool(redis_host, redis_port)
	n, err := removeKey(pool, ident)

	var status int
	var message string

	if err != nil {
		message = err.Error()
		status = 2
	} else if n > 0 {
		message = fmt.Sprintf("Host %v has been unregistered.", ident)
	} else {
		status = 1
		message = fmt.Sprintf("Host %v could not be found", ident)
	}

	fmt.Println(message)
	os.Exit(status)
}
