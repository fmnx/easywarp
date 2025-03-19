
# easywarp

🚀 **easywarp** is a lightweight tool that automatically registers and runs [Cloudflare WARP](https://developers.cloudflare.com/warp-client/) using the [WireGuard](https://www.wireguard.com/) protocol. It is designed to make using WARP fast and simple.

---

## ✨ Features

- Automatically registers and retrieves Cloudflare WARP configuration
- Runs WARP tunnel based on wireguard-go
- Lightweight and easy to deploy
- Cross-platform support (Linux / Windows)
- Manual route configuration for full traffic control

---

## 📦 Installation

```bash
git clone https://github.com/fmnx/easywarp.git
cd easywarp
go build
```

Alternatively, download a prebuilt binary from the [Releases](https://github.com/fmnx/easywarp/releases) page.

---

## 🚀 Usage

```bash
./easywarp
```

After starting, **easywarp** will:
1. Register and retrieve a Cloudflare WARP configuration.
2. Launch the WireGuard tunnel.

> 💡 **Note:** You need to manually configure system routes to direct traffic through the tunnel.

---

## ⚙️ Example Route Configuration (Linux)

To route all traffic through the WARP interface (`wg0`), use the following commands:

```bash
# Route all traffic through WARP interface (wg0)
ip route add 0.0.0.0/1 dev wg0
ip route add 128.0.0.0/1 dev wg0
```

You can also selectively route specific IP ranges through `wg0`.

---

## 💬 Contributing

Feel free to star ⭐, fork 🍴, or open an issue 🛠️ — contributions are welcome!

---

Let me know if you want a version with CLI flag descriptions, configuration examples, or systemd service setup — happy to expand it!