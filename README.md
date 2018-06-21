# RestAPI

### Dependencies:
```
go get github.com/gorilla/mux
go get github.com/go-sql-driver/mysql
```

### Build:

Set MySQL username and password inside main.go
```
a.Initialize("username", "password", "")
```

Build the program with
```
go build main.go app.go model.go
```
