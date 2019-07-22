Clone, get dependencies and running the application

go get -u github.com/gorilla/mux
go get -u github.com/gorilla/schema
cd $GOPATH/src; mv lenslocked lenslocked.com
cd $GOPATH/src/lenslocked.com
go run main.go