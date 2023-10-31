# pay-api
Payment API endpoint is created to simulate:
* Seeding a database with Faker data
* Pay API endpoint written in Go using
  * `go-pg` library
  * `pgx` library
* Generate a test workload with `k6`

### Generate Faker Data & Seed Database

* Run the seeding script for **Postgres**

```
pip install -r pay-api/requirements.txt
python db/seed.py --customers 100 --accounts 5 --payments 500 --db=postgresql://postgres:postgres@192.168.86.74:5432/bank
```
* Run the seeding script for **CockroachDB**
 * UI Console: `http://192.168.86.74:8080`
 * SQL Console: `postgresql://root@192.168.86.74:26257`

```
cockroach sql --insecure --host=192.168.86.74:26257
root@192.168.86.74:26257/defaultdb> CREATE DATABASE bank;
root@192.168.86.74:26257/defaultdb> \q

python db/seed.py --customers 100 --accounts 5 --payments 500 --db=postgresql://root@192.168.86.74:26257/bank
```

## Go Endpoints

*	POST: "/users", createUser
*	GET: "/users/:uuid/balances", checkBalances
*	POST: "/payments", sendPayment

### Execute Go: `gopg-api`

* Init the project

```shell
go mod init bun-api
go mod tidy
go get -u github.com/gin-gonic/gin
go get -u github.com/go-pg/pg/v10
```

* Start the API endpoint in http://localhost:8088

```
go run *.go
```

### Execute Go: `pgx-api`

* Init the project

```shell
go mod init pgx-api
go mod tidy

go get -u github.com/jackc/pgx/v4
go get github.com/gin-gonic/gin
go get github.com/jackc/pgx/v4/pgxpool@v4.18.1
```

* Start the API endpoint in http://localhost:8087

```
go run *.go
```

## K6 Stress Test

* `gopg-api` on port 8088

```
k6 run gopg-api/gopg-checkBalances.js
```

* `pgx-api` on port 8087

```
k6 run pgx-api/gopg-checkBalances.js
```

See sample results on this basic setup at [sample go-pg vs pgx results on the same workload](sample.log)
