# LogVoyage

Lets make new cool logging Saas :)

# Public consumers:
- HTTP consumer
- UDP consumer
- TCP consumer

All cosumers accepts send messages to one main consumers which sends data to Amazon SQS.
Then it will be consumed by workers and send to elastic / database, performed other metrics/queries.

# User
Every time we accept log we should increase size counters. Store it in redis and then sync to database?
- Need fast persistent storage!
- Traffic count by day
- Traffic limiting
- Store elastic url/index names

## Ideas
Cool site and logo animations
https://my.setapp.com/successful-registration


# Links
https://www.loggly.com/docs/syslog-ng-manual-configuration/

Awesome dashboard look
https://d13yacurqjgara.cloudfront.net/users/1113/screenshots/852453/attachments/90151/sentry-tv-dashboard-big.png
