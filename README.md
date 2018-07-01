# pgbeat
> Periodically insert current timestamp in your PostgreSQL database

Sometimes you need to measure replication lag between a primary and one or more standbys. When you use physical replication, you can use built-in functions like `pg_last_xact_replay_timestamp()` because system is always recovering. But, what if you are using logical replication? Maybe you are restoring a backup and you want to know the real backup timestamp. `pgbeat` comes to the rescue.

## Internals

`pgbeat` behavior is like a heartbeat system. It updates a given row at a given interval of time with the current timestamp. `pgbeat` works with replicas too. As soon as connected instance gets promoted, it will automatically start to update heartbeat.

## Highlights
* `pgbeat` uses an identifier (`-id`) associated with its timestamp so multiple daemons can run at the same time on the same instance without overlap.
* interval unit is seconds (`-interval`, `-recovery-interval`, `-timeout`). Milliseconds can be set using floating point value (ex: 0.25 for 250ms), except for `-timeout` where only integers are accepted.
* `pgbeat` relies on `libpq` for PostgreSQL connection. When `-host` is ommited, connection via unix socket is used. When `-user` is ommited, the unix user is used. And so on.
* `pgbeat` is able to create database with `-create-database` if it doesn't exist. A different database name can be specified with `-connect-database` to connect once and create the database (ex: "postgres" or "template1"). By default, libpq uses a database name based on the username to connect but this database might not exist so connection won't be successful.
* `pgbeat` is able to create table with `-create-table` if it doesn't exist.
* configuration file options **override** command-line arguments.

## Usage
Connect to a remote instance and prompt for password:
```
pgbeat -host 10.0.0.1 -port 5432 -user test -prompt-password -database test
```
Use a configuration file:
```
pgbeat -config config.yaml
```
Use both configuration file and command-line arguments:
```
pgbeat -config config.yaml -id 2 -interval 0.5
```
Print usage:
```
pgbeat -help
```

## License
`pgbeat` is released under [The Unlicense](LICENSE) license. Code is under public domain.
