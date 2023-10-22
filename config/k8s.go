//go:build k8s

package config

var Config = config{
	DB:    DBConfig{DSN: "root:root@tcp(mini-book-record-mysql:3308)/mini_ebook"},
	Redis: RedisConfig{Addr: "mini-book-record-redis:6380"},
}
