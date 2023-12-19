[![Test](https://github.com/snobb/go-dotenv/actions/workflows/test.yml/badge.svg)](https://github.com/snobb/go-dotenv/actions/workflows/test.yml)

# Go-DotEnv

A small implementation of nodejs dotenv library that populates env variables from a provided `.env` file (the file name is customisable).

This library is working particularly nice with Kelsey Hightower's envconfig.
https://github.com/kelseyhightower/envconfig

## Example Usage
The env file can have the following format:

```bash
$ cat .cenv
CONSUMER_PORT=8085
CONSUMER_LOG_LEVEL=debug
CONSUMER_CONCURRENCY=3
CONSUMER_KAFKA_TOPIC_NAME=data-topic
CONSUMER_KAFKA_GROUP_ID=cgroup
CONSUMER_KAFKA_BROKER_LIST=localhost:19092
```


With default name:
```Go
...

func init() {
    if err := dotenv.LoadEnv(); err != nil {
        slog.Error("dotenv", err)
    }
}

...
```

With a customised file name for env file.
```Go
...

func init() {
    if err := dotenv.LoadEnvFromFile(".custom_env"); err != nil {
        slog.Error("dotenv", err)
    }
}

...
```
