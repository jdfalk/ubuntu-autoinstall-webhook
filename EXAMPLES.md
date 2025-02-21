# Webhook API Examples

This document provides examples of how to call the webhook API.

## Sending a Webhook Request

### Example using cURL

```sh
curl -X POST "http://localhost:25000/api/webhook" \
     -H "Content-Type: application/json" \
     -d '{
       "origin": "curtin",
       "timestamp": 1440688425.6038516,
       "event_type": "start",
       "name": "cmd-install",
       "description": "curtin command install",
       "source_ip": "192.168.1.100"
     }'
```

### Example with Optional Fields

```sh
curl -X POST "http://localhost:25000/api/webhook" \
     -H "Content-Type: application/json" \
     -d '{
       "origin": "curtin",
       "timestamp": 1440688425.6038516,
       "event_type": "finish",
       "name": "cmd-install",
       "description": "curtin command install",
       "result": "SUCCESS",
       "files": [
         {
           "content": "fCBzZmRpc2s....gLS1uby1yZX",
           "path": "/var/log/curtin/install.log",
           "encoding": "base64"
         }
       ],
       "source_ip": "192.168.1.100",
       "status": "installing",
       "progress": 50,
       "message": "Installation is halfway done"
     }'
```

## Example using Python

```python
import requests
import json

url = "http://localhost:25000/api/webhook"
headers = {"Content-Type": "application/json"}
data = {
    "origin": "curtin",
    "timestamp": 1440688425.6038516,
    "event_type": "start",
    "name": "cmd-install",
    "description": "curtin command install",
    "source_ip": "192.168.1.100"
}
response = requests.post(url, headers=headers, data=json.dumps(data))
print(response.status_code, response.json())
```

## Example using Go

```go
package main

import (
 "bytes"
 "encoding/json"
 "fmt"
 "net/http"
)

func main() {
 url := "http://localhost:25000/api/webhook"
 data := map[string]interface{}{
  "origin": "curtin",
  "timestamp": 1440688425.6038516,
  "event_type": "start",
  "name": "cmd-install",
  "description": "curtin command install",
  "source_ip": "192.168.1.100",
 }

 jsonData, err := json.Marshal(data)
 if err != nil {
  fmt.Println("Error marshalling JSON:", err)
  return
 }

 resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
 if err != nil {
  fmt.Println("Error sending request:", err)
  return
 }
 defer resp.Body.Close()

 fmt.Println("Response Status:", resp.Status)
}
```

## Expected Response

```json
{
  "status": "success",
  "message": "Event received"
}
```
