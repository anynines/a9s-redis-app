# a9s Redis App

This is a sample app to check whether the a9s Redis service is working or not.

## Install, Push and Bind

Make sure you installed GO on your machine, [download this](https://golang.org/doc/install?download=go1.8.darwin-amd64.pkg) for mac.

Download the application
```
$ go get github.com/anynines/a9s-redis-app
$ cd $GOPATH/src/github.com/anynines/a9s-redis-app
```

Create a service on the [a9s PaaS](https://paas.anynines.com)
```
$ cf create-service a9s-redis40 redis-single-non-persistent-small myredis
```

Push the app
```
$ cf push --no-start
```

Bind the app
```
$ cf bind-service redis-app myredis
```

And start
```
$ cf start redis-app
```

At last check the created url...


## Local Test

Start Redis service with Docker:

```shell
$ docker run -d -p 6379:6379 redis redis-server --requirepass secret
```

Export a few environment variables and run the sample app:

```shell
$ export REDIS_HOST=localhost
$ export REDIS_PORT=6379
$ export REDIS_PASSWORD=secret
$ go build
$ ./a9s-redis-app
```

## Remark

To bind the app to other Redis services than `a9s-redis50`, have a look at the `VCAPServices` struct.
