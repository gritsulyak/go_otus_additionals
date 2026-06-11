module github.com/gritsulyak/go_otus_additionals/services/api

go 1.25.5

require (
	github.com/gritsulyak/go_otus_additionals/libs/logger v0.1.0
	github.com/lib/pq v1.10.9
)

replace github.com/gritsulyak/go_otus_additionals/libs/logger => ../../libs/logger
