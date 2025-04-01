# Cache System API Documentation

## API Endpoints

## 1. `/set` – Store a value in the cache

### Description:
The `/set` endpoint allows you to store a key-value pair in the cache with an expiration time (TTL).

### Request:
- **Method**: `POST`
- **Query Parameters**:
  - `key` (string) – The cache key.
  - `ttl` (integer, seconds) – The time-to-live (TTL) before the key expires.
- **Body**:
  - The value to store (can be a string, JSON, or any other data).

### Example `cURL` Request:
```bash
curl --location 'http://localhost:8080/set?key=myKey&ttl=300' --header 'Content-Type: application/json' --data 'This is my cached value'
```

### Expected Response:
*200 OK* - If the key is successfully stored.
*400 Bad Request* - If missing parameters


## 2. `/get` – Retrieve a value from the cache

### Description:
The `/get` endpoint allows you to retrieve a stored value from the cache using a key.

### Request:
- **Method**: `GET`
- **Query Parameters**:
  - `key` (string) – The cache key to retrieve.

### Example `cURL` Request:
```bash
curl --location 'http://localhost:8080/get?key=myKey'
```

### Expected Response:
*200 OK* - With the cached value in the response body.
*400 Bad Request* - If missing parameters
*404 Not Found* - If the key does not exist or has expired.

### Example Response:
```json
"This is my cached value"
```

## 3. `/trygetwithexpire` – Retrieve a value with expiration time

### Description:
The `/trygetwithexpire` endpoint retrieves a stored value from the cache along with its remaining time-to-live (TTL).

### Request:
- **Method**: `GET`
- **Query Parameters**:
  - `key` (string) – The cache key to retrieve.

### Example `cURL` Request:
```bash
curl --location 'http://localhost:8080/trygetwithexpire?key=myKey'
```
### Expected Response:
*200 OK* - If the key exists, returning the value and remaining TTL.
*400 Bad Request* - If missing parameters
*404 Not Found* - If the key does not exist or has expired.


### Example Response:
```json
{
  "value": "This is my cached value",
  "expires_in": "4m29s"
}
```

## 4. `/getKeys` – Retrieve values of specified keys

### Description:
The `/getKeys` endpoint allows you to retrieve the values of multiple keys at once.

### Request:
- **Method**: `POST`
- **Body**:
  - A JSON array containing the keys to retrieve.

### Example `cURL` Request:
```bash
curl --location 'http://localhost:8080/getKeys' --header 'Content-Type: application/json' --data '["myKey", "user123"]'
```
### Expected Response:
*200 OK* - With a JSON object containing the requested keys and their values.
*400 Bad Request* - If missing parameters
*500 Internal server error* - If the response cannot be serialized.

### Example Response:
```json
{
    "myKey": {
        "value": "This is my cached value",
        "expiration": "2025-03-27T09:53:29.3658301+01:00"
    },
    "user123": {
        "value": "Test user123",
        "expiration": "2025-03-27T09:53:38.2447624+01:00"
    }
}
```

## 5. `/list` – Retrieve all keys with truncated values and expiration times

### Description:
The `/list` endpoint returns all keys currently stored in the cache along with their values (truncated to 25 characters by default) and expiration times. If the optional `allValue` parameter is provided and set to `true`, the full values will be returned instead of truncated ones.

### Request:
- **Method**: `GET`
- **Query Parameters**:
  - `allValue` (optional, boolean) – If set to `true`, returns full values instead of truncated ones.

### Example `cURL` Request:
```bash
curl --location 'http://localhost:8080/list'
```

To retrieve full values:
```bash
curl --location 'http://localhost:8080/list?allValue=true'
```

### Expected Response:
*200 OK* – Returns a JSON array of all keys stored in the cache.

### Example Response (truncated values by default):
```json
[
    {
        "key": "myKey",
        "value": "This is my cached value...",
        "expires_in": "4m58.0297621s"
    }
]
```

### Example Response (full values with `allValue=true`):
```json
[
    {
        "key": "myKey",
        "value": "This is my full cached value with all content visible.",
        "expires_in": "4m58.0297621s"
    }
]
```

## 6. `/flush` – Clear the entire cache
### Description:
The `/flush` endpoint removes all keys and their associated values from the cache. This operation affects all nodes in the distributed system.

### Request:
- **Method**: `POST`

### Example `cURL` Request:
```bash
curl --location --request POST 'http://localhost:8080/flush'
```
### Expected Responses:
*200 OK*


## 7. `/remove` – Remove a key from the cache
### Description:
The `/remove` endpoint allows you to remove a specific key from the cache. If the key is not provided, the request will return a 400 Bad Request error.

### Request:
- **Method**: `GET`

- **Query Parameters**:
  - `key` (string) – The cache key to be removed.

### Example `cURL` Request:
```bash
curl --location --request POST 'http://localhost:8080/remove?key=myKey'
```
### Expected Responses:
*200 OK* - If the key was successfully removed from the cache.
*400 Bad Request* - If the key query parameter is missing.


## 8. `/removeallkeys` – Remove keys matching a pattern from the cache
### Description:
The `/removeallkeys` endpoint allows you to remove multiple keys from the cache that match a specified pattern. If the key parameter is missing, the request will return a 400 Bad Request error.

### Request:
- **Method**: `GET`

- **Query Parameters**:
  - `key` (string) – The pattern to match the keys for removal. This can be a part of the key or a wildcard pattern (e.g., user* to match all keys starting with "user").

### Example cURL Request:
```bash
curl --location --request DELETE 'http://localhost:8080/removeallkeys?key=user'
```

### Expected Responses:
*200 OK* - If keys matching the pattern were successfully removed from the cache.
*400 Bad Request* - If the key query parameter is missing.


## Internal Use Endpoints
These endpoints are used to synchronize the cache between all nodes in the system. 
They are for internal use only and are responsible for maintaining cache consistency across the network. 
While they can be called, they will be eventually be secured or modified to ensure that only internal services have access.

## 9. `/sync` – Synchronize cache between nodes
### Description:
The `/sync` endpoint is used internally to synchronize the cache between nodes. It ensures that cache data is consistent across all nodes in the network.

## 10. `/ping` – Ping to check node availability
### Description:
The `/ping` endpoint checks the availability of the node. It is used internally to monitor if the node is responsive.

## 11. `/export` – Export cache to another node
### Description:
The `/export` endpoint is used internally to export cache data to another node. This is part of the data synchronization between nodes.

## 12. `/diff` – Diff between local cache and another node's cache
### Description:
The `/diff` endpoint compares the local cache with another node's cache to identify the differences. It helps ensure data consistency across nodes.

## 13. `/set_batch` – Set multiple cache entries in a batch
### Description:
The `/set_batch` endpoint allows for setting multiple cache entries in a batch. This is used internally for bulk cache updates and synchronization.


# About config.json:

```json
{
    "port": ":8080",
    "num_counters": 10000000,
    "max_cost": 1073741824,
    "buffer_items": 64,
    "read_timeout": 5,
    "write_timeout": 5,
    "max_conns_per_ip": 1000,
    "max_requests_per_conn": 10,
    "peers": [
        "http://localhost:8080",
        "http://localhost:8081"
    ],
    "max_retries_to_disabled_node": 3,
    "heart_beat_interval_in_seconds": 5,
    "white_list_file_path": "whitelist.json"
}
```

## Configuration File Breakdown:
1. *port: ":8080"*
- This specifies the port where the server will listen for incoming requests. In this case, it's set to 8080.

2. *num_counters: 10000000*
- The number of counters (entries) the cache will manage. This controls how many items can be stored in the cache.

3. *max_cost: 1073741824*
- This is the maximum allowable cost for items in the cache. It is typically used in conjunction with memory management, where each item is assigned a "cost", and the total cost of all items cannot exceed this value.

4. *buffer_items: 64*
- This defines how many items the cache should buffer before writing to disk or synchronizing with other nodes. This helps optimize performance by reducing frequent I/O operations.

5. *read_timeout: 5*
- Specifies the timeout (in seconds) for reading requests. If a read operation takes longer than this time, it will be aborted.

6. *write_timeout: 5*
- Specifies the timeout (in seconds) for write operations. If writing data to the cache exceeds this time, it will be canceled.

*max_conns_per_ip: 1000*
- The maximum number of simultaneous connections allowed from a single IP address. This can help prevent abuse or overloading from a single source.

*max_requests_per_conn: 10*
- The maximum number of requests allowed per connection. This prevents a single connection from making too many requests, which could potentially overload the server.

*peers:*
- This is an array of other cache node addresses (peers) in the network. This is used for synchronizing data between nodes in a distributed cache setup.

*max_retries_to_disabled_node: 3*
- The number of retries a node will attempt to communicate with a peer before marking it as "disabled". This is useful for fault tolerance and managing communication failures.

*heart_beat_interval_in_seconds: 5*
- Defines the interval (in seconds) at which nodes send "heartbeat" signals to each other to check if the node is still active. If a node fails to respond, it will be marked as inactive and removed from the peer list.

*white_list_file_path: "whitelist.json"*
- This is the path to a JSON file containing the whitelist of allowed nodes. This file will be used to validate if incoming connections are from trusted peers.


## 🔒 About `whitelist.json`

The `whitelist.json` file defines a list of allowed nodes that can communicate with the distributed cache.  

### Structure:
```json
{
    "allowed_nodes": [
        "localhost:8080"
    ]
}
```

# Explanation:
**allowed_nodes**: An array of trusted nodes (host and port) that are permitted to interact with the cache.
- Any request from an unlisted node will be rejected.

**How to Use**:
- Add the IP or domain of authorized nodes inside the allowed_nodes array.

Example:

```json
{
    "allowed_nodes": [
        "localhost:8080",
        "192.168.1.10:8080",
        "node1.example.com:8081"
    ]
}
```