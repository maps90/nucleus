# Migration Package
migration package rely on config and storage package.

# How To Use
```go
import (
    "fmt"
    "github.com/maps90/nucleus/migration"
)

var DBName = "mydb"

func main() {
    // if status true then migration up
    var status = true

    version, _, err := migration.Migrate("resources/seeds/", DBName, status)
    if err != nil {
        return
    }
    fmt.Printf("DB Migration[%v] Completed \n", version)
}

```

# Create Migration File

__What You Need :__

- clone [migrate CLI](github.com/golang-migrate/migrate/cli).

- build ```go build -tags 'mysql' -o $GOBIN/migrate github.com/mattes/migrate/cli```


__Run the command :__
```
// this will create sql file under migrations folder.
#> migrate create -ext sql -dir migrations create_user
```

this will create up & down sql file, for more information [check out this link](github.com/golang-migrate/migrate) first.
