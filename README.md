# Dependencies
1. go get -u "github.com/google/jsonapi"
1. go get -u "github.com/gorilla/mux"
1. go get -u "golang.org/x/time/rate"

# To Run
go run main.go middleware.go

# Command Line Arguments
main.go [--host "localhost"] [--port 8080]

middleware.go [--max 10] [--burst 10]
