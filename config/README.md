# Config Package
config management using viper.

# How to use
```go

import (
    "fmt"

    "github.com/maps90/nucleus/config"
)

func main() {
    if err := config.New() {
        fmt.Println(err)
    }

    // custom filename
    // it will find security.yml in base folder.
    if err := config.New("security", "./") {
        fmt.Println(err)
    }
}

```

By default library will try to find `resources/application.yml` file.

# Configuration Example

Best Practice to write a configuration file.
```yaml
---
port: 3000
mysql:
    master:
        user: root
        password: root
        address: localhost:3306
        db: merchants
    slave:
        fallback_to: master
redis:
    master:
        host:
        user:
        password:
    slave:
        fallback_to: master

```