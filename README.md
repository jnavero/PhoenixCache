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

- 🔐 Security – Implement API Key authentication & optional SSL support.
- 🔐 Security – Protect internal endpoints


## 📖 API Documentation

For detailed information about the configuration and available API endpoints, please check the [API Documentation](API_Documentation.md).