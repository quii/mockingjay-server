# mockingjay server

Mockingjay lets you define the contract between a consumer and producer with a simple configuration file. 

Mockingjay will then give you fake services to write integration tests against and [consumer driven contracts](http://martinfowler.com/articles/consumerDrivenContracts.html) to run against your real downstream services.

Mockingjay is fast, requires no coding and is better than other solutions because it will ensure your mock servers and real integration points are consistent.

## Rationale

In the hip exciting world of SOA/microservices with heavy investment in PaaS/IaaS you want to be able to quickly iterate over small services and deploy to live quickly and without fear of breaking things. 

You will probably employ things like versioning to help but you might also be spending time writing consumer driven contracts (CDCs) to ensure your integration points are working. 

In addition you might be writing integration tests against fakes/stubs to ensure your code can send the correct requests and be able to parse responses. 

![alt tag](http://i.imgur.com/oC6BjGn.png)

If you squint hard enough, you can imagine that the requirements for both CDCs and fake servers are the same. *Given a particular request, I expect a particular kind of response*. Yet with this set up you are duplicating this information with code in two different files which obviously isn't ideal. 

What mockingjay enables you to do is to capture these requirements in one configuration file. 

````yaml
---
 - name: My very important integration point
   request:
     uri: /hello
     method: GET
   response:
     code: 200
     body: '{"message": "hello, world"}'
     headers:
       content-type: application/json

# define as many as you need...
````

From this you can create a fake server to write integration tests with and also check the service you are dependant on is consistent with what you expect. 

#### Main advantages

- No coding whatsoever, so no naughtiness in fake servers overcomplicating things. Even non developers can add new scenarios to test with.
- The contract is defined once, rather than dispersed across different scripts which you have to keep in sync. 
- Entirely language agnostic. If you speak HTTP you can use mockingjay.
- Checks the structure of the data (currently JSON is the only type checked) rather than the contents, which will reduce flakiness of your builds.
- Efforts have been made to make sure mockingjay is fast, to help you keep your feedback loops tight. 

#### Drawbacks/constraints

- You can only express your consumer-producer interaction in terms of isolated request/responses. Sometimes you might need to test a number of requests which are dependant on each other; but hopefully this scenario is rare.

## Installation

     $ go get github.com/quii/mockingjay-server

## Running a fake server

````bash
$ mockingjay-server -config=example.yaml -port=1234 &
$ 2015/04/13 14:27:54 Serving 3 endpoints defined from example.yaml on port 1234
$ curl http://localhost:1234/hello
$ {"message": "hello, world"}
````

## Check configuration is compatible with a real server

````bash
$ mockingjay-server -config=example.yaml -realURL=http://some-real-api.com
$ 2015/04/13 21:06:06 Test endpoint (GET /hello) is incompatible with http://some-real-api - Couldn't reach real server
$ 2015/04/13 21:06:06 Test endpoint 2 (DELETE /world) is incompatible with http://some-real-api - Couldn't reach real server
$ 2015/04/13 21:06:06 Failing endpoint (POST /card) is incompatible with http://some-real-api - Couldn't reach real server
$ 2015/04/13 21:06:06 At least one endpoint was incompatible with the real URL supplied
````

## Building

### Requirements

- Go 1.3+ installed ($GOPATH set, et al)
- godep https://github.com/tools/godep
- golint https://github.com/golang/lint


    ./build.sh

### TODO

- Pretty diagrams explaining it all (interations between fakes, CDCs et al.)
- Currently the CDC part only checks status codes. 
- Check XML structure when applicable
- Get the binaries built somewhere for non gophers
- Investigate a more standard test output (http://windyroad.com.au/2011/02/07/apache-ant-junit-xml-schema/)

### Things to figure out

- Still somewhat reliaint on golden data. i.e fake /user/2 - does user 2 exist on the real server? What can be done?
- Is it possible to do chains of requests for more complicated tests but still keep it nice and simple. Should we?
