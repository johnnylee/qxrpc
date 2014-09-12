qxrpc
=====

RPC between a Golang HTTP server and a Qooxdoo Javascript application.

The code consists of a small Javascript class and help functions in Go. 

On the Go side, a basic HTTP server handler function:
```Go
func add(w http.ResponseWriter, r *http.Request) {
	args := struct{ X, Y int }{}
	if !qxrpc.ParseArgs(w, r, 1024, &args) { return }

	// Compute the result.
	qxrpc.SendResponse(w, args.X + args.Y)
}
```

On the Javascript side: 
```Javascript
var client = new qxrpc.Client("http://localhost:8000");
  client.callAsync(
    // success is boolean, error is 
    function(success, code, message, result)
    {
      if (success) {
        alert(result);
      } else {
        alert(message);
      }
    },
    "/rpc-1.0/add", 
    { X: 1, Y: 2 });
});
```
