go test ./... -coverprofile=coverage.out
CALL coverage_badge.bat
go tool cover -html=coverage.out
