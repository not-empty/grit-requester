reset
go test -v -tags=unit -coverpkg=./... -coverprofile=./coverage/coverage-unit.out ./
go tool cover -func=./coverage/coverage-unit.out -o ./coverage/coverage-unit.txt
go tool cover -html=./coverage/coverage-unit.out -o ./coverage/coverage-unit.html