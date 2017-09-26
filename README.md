=Cockroach BugReport 1=

This repo demonstrates a bug in CockroachDB. To execute, run:

```
cockroach start --insecure &
cockroach user set maxroach --insecure
cockroach sql --insecure -e 'CREATE DATABASE bank'
cockroach sql --insecure -e 'GRANT ALL ON DATABASE bank TO maxroach'
```

and run

```
go run ./main.go
```

Expected result: No insertion failures

Actual result: INSERTs inside a transaction fail 0.5% - 1.0% of the time