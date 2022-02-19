module shortly.db.clients

go 1.17

require (
	github.com/go-redis/redis/v8 v8.11.4
	github.com/stretchr/testify v1.5.1
	shortly.data.data_model v0.0.1
)

replace shortly.data.data_model => ../data

require (
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)
