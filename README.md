# server-monitoring - Monitor your servers

Status Page for monitoring your websites and applications with graphs and analytics.

## Run everywhere

`server-monitoring` is written in Go(lang). All you need is the precompiled binary based on your operating system, and the
HTML/CSS/Javascript files. You can even run `server-monitoring` on your Raspberry Pi.

## Notifications

`server-monitoring` can notify you via email or Telegram when a service is unavailable. More notification types coming soon.

It is possible to use your own template for notifications. The [golang template engine](https://golang.org/pkg/text/template/#example_Template) is used for this purpose. Possible variables are `{{.Name}}`, `{{.Reason}}` and `{{.Date}}`.

## Run on Docker

#Use the [official Docker image](https://hub.docker.com/r/koloooo/monhttp) to run monhttp in seconds.

``` shell
#
```

To save the config.env from the container for later, you need to mount the path `/monhttp/config`. Add a volume for this
when starting the container.

``` shell
#
```
## Use mongodb
Use this tutorial to install MongoDB 5.0 Community Edition on LTS (long-term support) releases of Ubuntu Linux using the apt package manager.
 use this link [`https://docs.mongodb.com/manual/installation/`](https://docs.mongodb.com/manual/installation/) for installation in your server.

## Use docker-compose

Simply run `docker-compose up` to start `monhttp` together with a postgres database. Open
up [`http://localhost:8081`](http://localhost:8081) in your browser and enjoy `monhttp`. The default user is `admin` and
the password is `admin` too.

## Build it locally

Make sure you have Go 1.15 installed on your computer. Clone the repository and execute the build
command.

``` shell
git clone git@github.com:irmitcod/server-monitoring.git
cd server-monitoring
make buildLocal
```

Then you will find all files and folders in the dist folder. Change to the dist folder and start monhttp
with `./server-monitoring`.

## Configuration

After the initial setup, there is a config.json in the config folder. This file can be used to change or save the
configuration.

| Key  | Value  | Description  |
|---|---|---|
| DATABASE_HOST | localhost  |   |
| DATABASE_NAME |  monitoring |   |
| DATABASE_PASSWORD |  |   |
| DATABASE_PORT |   |   |
| DATABASE_USER |   |   |
|   |   |   |
| NOTIFIER |   |   |
|   |   |   |
| SERVER_PORT | 8081  |   |
|   |   |   |
| USERS | admin:123456,  | A list in the format "name:password" you can add here as many users as you want to  |
|   |   |   |
| SCHEDULER_ENABLED  | true  | If false, then no data is collected  |
| SCHEDULER_NUMBER_OF_WORKERS  | 5  | How many "workers" should process the services asynchronously. If there are many services, the value should be increased.  |


You can also use environment variables to configure `server-monitoring`. Environment variables override the values from the `config.json` file.
