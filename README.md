# go-snippets

Nout:go-snippets vvv$ go run http_post_json.go 

Nout:go-snippets vvv$  ab -c 20 -n 10000 -p postfile.json -T 'application/json' http://127.0.0.1:8008/
Nout:go-snippets vvv$ curl -i -X POST -H "Content-Type: application/json" -d '{"url": "google.com", "id": "05461bd4-f3b7-46c7-92ce-9c5fdc662b47"}' http://127.0.0.1:8008/
HTTP/1.1 200 OK
