```bash
go install github.com/vektra/mockery/v2@latest

go get github.com/stretchr/testify/mock@v1.10.0

mockery --all --output=./mocks --outpkg=mocks --recursive


# test:
go test ./internal/api/system -v


```