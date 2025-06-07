# Go Reverse Proxy Load Balancer

A simple round-robin load balancer written in Go, using the `httputil.ReverseProxy` from the standard library. It forwards incoming requests to a list of backend servers.

---

## 🗂️ Project Structure

```
go-load-balancer/
├── go.mod
├── README.md
└── src/
└── main.go

```

---

## 🚀 How It Works

- The load balancer listens on a specified port (default: `:8000`).
- Incoming HTTP requests are distributed in round-robin fashion to the list of backend servers (`google.com`, `wikipedia.org`, `github.com`).
- Currently, each server is considered *alive* (always returns `true` in `isAlive()`).
- Utilizes Go’s built-in `httputil.ReverseProxy`.

---

## 📦 Getting Started

### Prerequisites

- Go 1.18+

### Clone the Repository

```bash
git clone https://github.com/thatquietkid/go-load-balancer.git
cd go-load-balancer
```

### Run the Project

```bash
cd src
go run main.go
````

Then visit:

```
http://localhost:8000
```

Each refresh forwards the request to the next backend server.

---

## 🔧 Configuration

To add or change the backend servers, modify this block in `main.go`:

```go
servers := []Server{
    newSimpleServer("https://www.google.com"),
    newSimpleServer("https://www.wikipedia.org"),
    newSimpleServer("https://www.github.com"),
}
```

To change the listening port, update:

```go
lb := newLoadBalancer(":8000", servers)
```

---

## ✅ To-Do

* [ ] Implement actual health checks in `isAlive()`
* [ ] Add support for sticky sessions
* [ ] Add weighted round-robin mechanism
* [ ] Improve file structure (e.g., separate `server.go`, `balancer.go`, `main.go`)
* [ ] Graceful shutdown and logging enhancements
* [ ] Dockerize the project for deployment

---

## 📄 License

This project is open source and available under the [MIT License](LICENSE).
