# LogVoyage

Open Source Log, Exception, Metrics management.

This project is under heavy development. Any contributions would be appreciated. Before starting a feature or fix please create pull request. [Join our chat](http://link) or write me a message.

## Running
Preferred way for running is Docker.  Start dependent services:
```
start postgresql
start elasticsearch
start rabbitmq
```
Start LogVoyage services
```
start producer
start consumer
start logvoyage
start frontend
```

## Development
### Configuration
LogVoyage uses config file in json format. Config should be placed in $HOME directory `~/.logvoyage/config.json`.
Application configuration can be overridden via `LV` prefix. For example nested json key `db.database`
can be changed using `LV_DB_DATABASE=dbname` env variable. See example config file `config/config.json`.

### Installation
```
go get -u bitbucket.org/firstrow/logvoyage
```
