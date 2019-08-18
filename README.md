Clone, go get dependencies and running the application

go get -u github.com/gorilla/mux
go get -u github.com/gorilla/schema
go get -u github.com/lib/pq
go get -u github.com/jinzhu/gorm
go get -u golang.org/x/crypto/bcrypt

cd $GOPATH/src; mv lenslocked lenslocked.com
cd $GOPATH/src/lenslocked.com
go run main.go


#------ postgres -----
Install: https://postgresapp.com/
cmd> psql -U postgres
psql> create database lenslocked_dev; // first time only
psql> \c lenslocked_dev
psql> select * from users;