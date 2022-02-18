module shortly.workers

require (
	github.com/go-redis/redis/v8 v8.11.4
	shortly.data.data_model v0.0.1
	shortly.db.clients v0.0.1
	shortner v0.0.1
)

require (
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
)

replace (
	shortly.data.data_model => ../data
	shortly.db.clients => ../db_clients
	shortner => ../shortner_mod
)

go 1.17
