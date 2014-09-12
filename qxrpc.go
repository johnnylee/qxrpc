package qxrpc

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

// Error codes returned internally.
const (
	CodeRequestTooLarge = 131
	CodeRead            = 132
	CodeJSONParse       = 133
)

// ParseArgs parses arguments into an interface. It returns true on
// succcess, and false otherwise.
//
// Usage in a hanlder might look like this:
//
// args := struct{ X, Y int }{}
// if !ParseArgs(w, r, 1024, &args) { return }
func ParseArgs(
	w http.ResponseWriter, r *http.Request,
	maxSize int64, args interface{}) bool {

	// Check the request content-length.
	if r.ContentLength > maxSize {
		SendError(w, CodeRequestTooLarge, "Request too large")
		return false
	}

	// Read in the body.
	buf := make([]byte, r.ContentLength)
	_, err := io.ReadFull(r.Body, buf)
	if err != nil {
		SendError(w, CodeRead, err.Error())
		return false
	}

	// Unmarshal the body into args.
	if err = json.Unmarshal(buf, args); err != nil {
		SendError(w, CodeJSONParse, err.Error())
		return false
	}

	return true
}

// SendError sends the given error code and message to the caller.
// Error codes from 100-200 are reserved.
func SendError(w http.ResponseWriter, errCode int64, errMsg string) {
	sendQx(w, map[string]interface{}{
		"success": false,
		"code":    errCode,
		"message": errMsg,
	})
}

// SendResponse sends a successful response to the caller. The result
// parameter must be serializable as json.
func SendResponse(w http.ResponseWriter, result interface{}) {
	sendQx(w, map[string]interface{}{
		"success": true,
		"result":  result,
	})
}

func sendQx(w http.ResponseWriter, resp map[string]interface{}) {
	w.Header().Set("Content-Type", "application/json")
	b, err := json.Marshal(resp)
	if err != nil {
		log.Printf("Failed to marshal response:\n    %v\n", err)
		return
	}

	if _, err = w.Write(b); err != nil {
		log.Printf("Failed to write response:\n    %v\n", err)
		return
	}
}
