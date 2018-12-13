# MYSQL
package mysql source designed to give lazy load singleton access to mysql connections
it doesn't provide any cluster nor balancing suport, assuming it is handled
in lower level infra, i.e. proxy, cluster etc.

# Configuration Example
configuration can be in yaml or json, example :

```
mysql:
  log: true
  max_open: 30
  max_lifetime: 30
  account:
    user: root
    password: admin123
    address: localhost:3306
    db: lumo_account
```