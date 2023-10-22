//go:build !k8s

package config

var Config = config{
	DB:    DBConfig{DSN: "root:root@tcp(localhost:13316)/mini_ebook"},
	Redis: RedisConfig{Addr: "localhost:6379"},
}
