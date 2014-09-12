

qx.Class.define(
  "qxrpc.Client",
  {
    extend: qx.core.Object,
    
    construct: function(address) {
      this.address = address
    },

    statics: {
      errorCodes:
      {
        abort           = 101,
        error           = 102,
        timeout         = 103,
        requestTooLarge = 131,
        readError       = 132,
        parseError      = 133
      }
    },
    
    members:
    {
      callAsync: function(handler, path, args) 
      {
        this.debug("callAsync: calling " + this.address + path)
        var req = new qx.io.request.Xhr(this.address + path, "POST");
        req.setAsync(true);
        req.setCache(false);
        req.setRequestData(qx.lang.Json.stringify(args));
        
        req.addListener("success", function(e) {
          var resp = e.getTarget().getResponse();
          handler(resp.success, resp.code, resp.message, resp.result);
        }, this);
        
        req.addListener("abort", function(e) {
          handler(false, this.abort, "Request aborted.", null);
        }, this);
        
        req.addListener("error", function(e) {
          handler(false, this.error, "Unknown error.", null);
        }, this);
        
        req.addListener("timeout", function(e) {
          handler(false, this.timeout, "Request timed out.", null);
        }, this);
        
        req.send();
      }
    }
  });

