# monk-scheduller

Required env
- DB_HOST
- DB_PORT
- DB_NAME
- DB_USER
- DB_PASS
- APP_PORT
- APP_HOST

Required soft
- [flyway](https://flywaydb.org/documentation/commandline/#download-and-installation)

```
#!/bin/bash

go run main.go setup

flyway migrate

go build main.go

./main start
```