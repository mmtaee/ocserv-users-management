# Mockery installation
```bash
go install github.com/vektra/mockery/v2@latest

go get github.com/stretchr/testify/mock@v1.10.0

mockery --all --output=./mocks --outpkg=mocks --recursive
```

# commands
```bash
# test:
go test ./...

# Start Server
go run main.go serve --debug

# List of Admin users in cli
go run main.go admins
 
# Create Admin User in cli
go run main.go create-admin -u USERNAME -p PASSWORD
```