# PhoenixCache – A Distributed In-Memory Cache

**PhoenixCache** is a high-performance, distributed in-memory cache designed for speed, resilience, and scalability. 
It ensures data consistency across multiple nodes with efficient synchronization and automatic recovery mechanisms.

# 🔥 Features
- ⚡ Fast & Lightweight – Optimized for low-latency caching.
- 📡 Distributed & Synchronized – Multi-node support with push-based updates.
- 💾 Auto-Recovery – Nodes can recover missing data upon reconnection.
- 📡 Peer Monitoring – Heartbeat mechanism to detect active nodes.
- ♻️ Expiry-Based Cleanup – No need for periodic cache sweeps.

# 📖 How It Works
- Nodes synchronize via HTTP when data is modified (Set, Remove, Flush).
- If a node goes offline, it will automatically catch up when it reconnects.
- A diffing mechanism ensures stale data is refreshed on reactivation.
- Configurable whitelist of allowed peers for security.

# 🛠️ Getting Started
Clone the repo
Configure nodes & authentication
Run multiple instances
Enjoy blazing-fast distributed caching! 🚀


# ✅ Next Steps

🔐 Security – Implement API Key authentication & optional SSL support.


# About config.json File:
Here is an explanation of the configuration file for your distributed memory cache project:
```
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

How to use:
Modify the port if you wish to run the server on a different port.
Update the peers array to add the addresses of any other nodes that should be synchronized.
Adjust the timeouts, retry settings, and connection limits to suit your application's needs and expected load.
The whitelist file path is used to manage security by allowing only trusted nodes to connect to your cache network.

