### Messagebird Demo app

Requirements
------------
- [Sign up](https://www.messagebird.com/en/signup) for a free MessageBird account
- Create a new access key in the developers sections

Installation
------------
The easiest way to use this demo MessageBird API app in your Go project is to install it using *go get*:

```
$ go get github.com/akhilesharora/hire
```

Examples
--------
Here is a quick example on how to get started. Assuming the **go get** installation worked, you can import the messagebird package like this:

Then, to create an instance of *messagebird.Client*, add Messagebird Access key to *config.default.json*
```
{
  "Version": "dev",
  "ServerAddr": "localhost:8070",
  "LogLevel": "debug",
  "AccessKey": "ACCESS_KEY",
  "Originator": "Messagebird"
}
```

Once all dependencies are resolved and then start the server.
```sh
$ go build
$ ./hire server
```

Example cURL request

```sh
$ curl -X POST   http://localhos8070/  -d '{"recipient":YOUR_PHONE_NUMBER,"originator":"Messagebird","message":"This is a test message."}'
```

