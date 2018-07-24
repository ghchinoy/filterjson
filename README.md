# filterjson

a small service to remove top level keys from a JSON structure


## usage

API is `POST /` accepting a body of `application/json` and query params `remove` and/or `filter`

* `filter` will keep only the field name given in `filter` (1..*)
* `remove` will remove the field name given in `remove` (1..*)

Accepts either a JSON object or an array of JSON objects (but not, say, an array 
of arrays)

examples with `remove`

```
curl -v -X POST -d @samples/array.json -H "content-type: application/json" 'http://localhost:12001?remove=name'
```

```
curl -v -X POST -d @samples/geo.json -H "content-type: application/json" 'http://localhost:12001?remove=geometry'
```

Adds in HTTP headers `Removes`, an array of remove query params and 
`Unfiltered-Content-Length`, the original content length

example:
```
Removes: [geometry]
Unfiltered-Content-Length: 3760
```

## build / run

For development, can just `go run main.go` from the `src` dir

To build a container, `docker build -t ghchinoy/filterjson .` using the 2-stage 
build Dockerfile, then to run it `docker run -d -p 8080:12001 --name filterjson ghchinoy/filterjson` (this exposes the internal `12001` 
to make the service available on `8080`)

There's also a nodejs version of the server in `nodejs` that implements only the 
`filter` param (not `remove`), and for top level JSON objects (not JSON arrays)
