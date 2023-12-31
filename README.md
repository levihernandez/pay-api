# pay-api
Payment API endpoint is created to simulate:
* Seeding a database with Faker data
* Pay API endpoint written in Go using
  * `go-pg` library
  * `pgx` library
* Generate a test workload with `k6`

### Generate Faker Data & Seed Database

* Install dependencies for Postgres and CockroachDB

```shell
pip install -r pay-api/requirements.txt
```

* Run the seeding script for **Postgres**

```
python db/seed.py --customers 100 --accounts 5 --payments 500 --db=postgresql://postgres:postgres@192.168.86.74:5432/bank
```
* Run the seeding script for **CockroachDB**
 * UI Console: `http://192.168.86.74:8080`
 * SQL Console: `postgresql://root@192.168.86.74:26257`

```
cockroach sql --insecure --host=192.168.86.74:26257
root@192.168.86.74:26257/defaultdb> CREATE DATABASE bank;
root@192.168.86.74:26257/defaultdb> \q

python db/seed.py --customers 100 --accounts 5 --payments 500 --db=cockroachdb://root@192.168.86.74:26257/bank
```

## Go Endpoints

*	POST: "**/users**", createUser
*	GET: "**/users/:uuid/balances**", checkBalances
*	POST: "**/payments**", sendPayment

### Execute Go: `gopg-api`

* Init the project

```shell
cd gopg-api/
go mod init gopg-api
go mod tidy
```

* Start the API endpoint in http://localhost:8088

```
go run *.go
```

### Execute Go: `pgx-api`

* Init the project

```shell
cd pgx-api/
go mod init pgx-api
go mod tidy
```

* Start the API endpoint in http://localhost:8087

```
go run *.go
```

## K6 Stress Test

Modify your endpoint ports in each project `config.json` to point to CockroachDB.

* `gopg-api` on port 8088

```
k6 run gopg-api/gopg-checkBalances.js
```

* `pgx-api` on port 8087

```
k6 run pgx-api/gopg-checkBalances.js
```

See sample results on this basic setup at [sample go-pg vs pgx results on the same workload](sample.log)
