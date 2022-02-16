module shortner

go 1.17

require (
	github.com/btcsuite/btcutil v1.0.2
	shortly.data.data_model v0.0.1
	)

replace (
	shortly.data.data_model => ../data
)
