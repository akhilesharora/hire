### Installation

Install the dependencies and devDependencies and start the server.

```sh
$ go build
$ ./hire server
```
Example cURL request

```sh
$ curl -X POST   http://localhos8070/   -H 'cache-control: no-cache'  -d '{"recipient":YOUR_PHONE_NUMBER,"originator":"Messagebird","message":"This is a test message."}'
```
