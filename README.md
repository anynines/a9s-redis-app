# a9s Redis App

This is a sample app to check whether the Redis service is working or not.

## Install, push and bind

Download the application
```
$ go get github.com/anynines/a9s-redis-app
$ cd $GOPATH/src/github.com/anynines/a9s-redis-app
```

Create a service on the [a9s PaaS](https://paas.anynines.com)
```
$ cf create-service a9s-redis32 redis-single-non-persistent-small myredis
```

Push the app
```
$ cf push --no-start
```

Bind the app
```
$ cf bind-service redis-app myredis
```

And restage
```
$ cf restage redis-app
```

At last check the created url...


## Local test

To start it locally you have to export the env variable VCAP_SERVICES
```
$ export VCAP_SERVICES='{
  "a9s-redis32": [
   {
    "credentials": {
     "host": "localhost",
     "password": "secret",
     "port": 6379
    }
   }
  ]
 }'
 ```

Start Redis service with Docker:
```shell
$ docker run -d -p 6379:6379 redis redis-server --requirepass secret
```

Run the sample app
```
$ go build
$ ./a9s-redis-app
```

## Remark

To bind the app to other Redis services than `a9s-redis32`, have a look at the `VCAPServices` struct.
