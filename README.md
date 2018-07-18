# filterjson

a small service to remove top level keys from a JSON structure

API is `POST /` with required 1..* query param `remove`

Will remove the field name given in `remove`

Accepts both a top-level JSON array object or JSON object

```
curl -v -X POST -d @samples/array.json -H "content-type: application/json" 'http://localhost:12001?remove=name'
```

```
curl -v -X POST -d @samples/geo.json -H "content-type: application/json" 'http://localhost:12001?remove=geometry'
```

Adds in HTTP headers `Removes`, an array of remove query params and `Unfiltered-Content-Length`, the original content length

example:
```
Removes: [geometry]
Unfiltered-Content-Length: 3760
```