# pgbeat
> Periodically insert current timestamp in your PostgreSQL database

Sometimes you need to measure replication lag between a primary and one or more standbys. When you use physical replication, you can use built-in functions like `pg_last_xact_replay_timestamp()` because system is always recovering. But, what if you are using logical replication? Maybe you are restoring a backup and you want to know the real backup timestamp. `pgbeat` comes to the rescue.

## Internals

`pgbeat` behavior is like a heartbeat system. It updates a given row at a given interval of time with the current timestamp.

## Highlights
`pgbeat` uses an identifier (`-id`) associated with its timestamp so multiple daemons can run at the same time on the same instance without overlap.

Interval unit is milliseconds (`-interval`).

`pgbeat` relies on `libpq` for PostgreSQL connection:
* when `-host` is ommited, connection via unix socket is used
* when `-user` is ommited, the unix user is used
* default `-database` is `postgres`
* `-password` is optional

You will have to use CTRL+C (SIGINT) or kill (SIGTERM) to terminate `pgbeat`.

Configuration file options **override** command-line arguments.

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
pgbeat -config config.yaml -id 2 -interval 500
```
Print usage:
```
pgbeat -help
```

## License
`pgbeat` is released under [The Unlicense](https://github.com/jouir/pgbeat/blob/master/LICENSE) license. Code is under public domain.
