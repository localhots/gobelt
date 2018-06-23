i:
	@go install ./...

gen:
	@go run set/internal/gen/main.go -tpl=set/internal/impl -dest=set

test:
	@go test ./...

dbtest:
	@go test ./sqldb -dsn="root:@(127.0.0.1:3306)/sqldb_pkg_test"

fulltest: test dbtest