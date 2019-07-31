Clone, go get dependencies and running the application

go get -u github.com/gorilla/mux
go get -u github.com/gorilla/schema
go get -u github.com/jinzhu/gorm
cd $GOPATH/src; mv lenslocked lenslocked.com
cd $GOPATH/src/lenslocked.com
go run main.go


#------ postgres -----
cmd> psql -U postgres
psql> \c lenslocked_dev
psql> select * from users;