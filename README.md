# topgg-server
a webserver written in golang to that recieves vote webhooks from Top.gg then pipes them into RabbitMQ  

## usage
```sh
go mod download
go build main.go
# see env's below
./main
```

### environment variables
| variable   | description                               | default            |
|------------|-------------------------------------------|--------------------|
| ENDPOINT   | the route vote data will be sent to       | /webhooks/votes    |
| ADDRESS    | the address (:port) to run the server at  | :4500              |
| DBL_SECRET | the secret                                |                    |
| AMQP_URL   | the URL the amqp server is running behind | amqp://localhost// |
| AMQP_GROUP | the amqp group to publish data to         | votes              |