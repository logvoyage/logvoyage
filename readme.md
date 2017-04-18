# LogVoyage
Open Source Log, Exception, Metrics management.

## Installation
```
go get -u bitbucket.org/firstrow/logvoyage
```

## Configuration
LogVoyage uses config file in json format. Config should be placed in $HOME directory `~/.logvoyage/config.json`.
Application configuration can be overridden via `LV` prefix. For example nested json key `db.database`
can be changed using `LV_DB_DATABASE=dbname` env variable. See example config file `config/config.json`.

# TODO:
https://www.balabit.com/sites/default/files/documents/syslog-ng-ose-latest-guides/en/syslog-ng-ose-guide-admin/html/loggen.1.html
- Test using loggen

# Roadmap v1
- docs
- docker images
- fully working logs