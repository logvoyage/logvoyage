# LogVoyage

## Installation
```
go get -u bitbucket.org/firstrow/logvoyage
```

## Configuration
LogVoyage uses config file in json format. Config should be placed in $HOME directory `~/.logvoyage/config.json`.
Application configuration can be overridden via `LV` prefix. For example nested json key `db.database`
can be changed using `LV_DB_DATABASE=dbname` env variable. See example config file `config/config.json`.