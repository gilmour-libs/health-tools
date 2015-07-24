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

func errExit(err error) {
	fmt.Println(err)
	os.Exit(2)
}

func newRedisPool(host string, port int) *redis.Pool {
	return backends.GetPool(fmt.Sprintf("%v:%v", host, port))
}

func removeKey(pool *redis.Pool, ident string) {
	if ident == BLANK {
		err := errors.New("Identifier to remove cannot be blank")
		errExit(err)
	}

	conn := pool.Get()
	defer conn.Close()

	n, err := redis.Int(conn.Do("HDEL", GilmourHealthKey, ident))
	if err != nil {
		errExit(err)
	}

	if n > 0 {
		fmt.Printf("%v has been unsubscribed from health checks.\n", ident)
	} else {
		fmt.Printf("%v is not a registered host entry\n", ident)
	}
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
	removeKey(pool, ident)
	os.Exit(0) //Exit cleanly
}
