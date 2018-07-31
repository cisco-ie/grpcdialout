# mdt-go

To get going...

1. From in this directory:

    ```
    export GOPATH=`pwd`
    ```

1. Run `gen_from_protos.sh`, which will install GRPC tools and library dependencies, and compile protofiles.
2. Build and install the server using `go install server`
3. Run the server using `./bin/server`

Once you have data being sent to port 2345 wherever you are running the server, success looks like:

```
$ ./bin/server
Hello, world.
Receiving dialout stream from 10.55.106.9:60439!
ReqId = 0!
ReqId = 1!
ReqId = 2!
ReqId = 3!
ReqId = 4!
ReqId = 5!
ReqId = 6!
ReqId = 7!
ReqId = 8!
ReqId = 9!
ReqId = 10!
...etc...
```

Really not very interesting yet!
