# NATS Fundamentals Workshop, Part I

## Configure NATS CLI

Create a context to connect to the remote NATS server

	nats context add rethink --server nats://52.53.177.42:4222 --select

Verify the context is created and selected 

	nats context ls

Check the connection to the server

	nats server check connection

## Publish-Subscribe

	nats sub SUBJECT.NAME.STRING

	nats pub SUBJECT.NAME.STRING "PAYLOAD"

## Publish-Subscribe with Wildcard

	nats sub "SUBJECT.*.STRING"

	nats pub SUBJECT.NAME.STRING "PAYLOAD"

## Fan Out

	nats sub SUBJECT.NAME.STRING

	nats sub "SUBJECT.*.STRING"

	nats sub "SUBJECT.NAME.*"

	nats pub SUBJECT.NAME.STRING "PAYLOAD"

## Queue Group

	nats sub SUBJECT.NAME.STRING --queue UNIQUE_GROUP_NAME

	nats sub SUBJECT.NAME.STRING --queue UNIQUE_GROUP_NAME

	nats sub SUBJECT.NAME.STRING --queue UNIQUE_GROUP_NAME

	nats pub SUBJECT.NAME.STRING "Payload"

## Request-Reply with Header

	nats reply SUBJECT.NAME.STRING "Response message"

	nats request SUBJECT.NAME.STRING "Request message" --header Key:Value

## NATS Service

	nats service serve rethink

	nats request rethink.echo "Hello, NATS!" 

## Create a Stream

	nats stream add STREAM_NAME --subjects SUBJECT.NAME.STRING --defaults

## Publish a Message to the Stream

	nats pub SUBJECT.NAME.STRING "PAYLOAD"

	nats stream info STREAM_NAME 

## Create a Consumer

	nats consumer add STREAM_NAME CONSUMER_NAME --defaults

	nats consumer info STREAM_NAME CONSUMER_NAME

## Consume the Message

	nats consumer next STREAM_NAME CONSUMER_NAME

	nats stream info STREAM_NAME 

