# Build A Loadbalancer with Golang

---
## Introduction
This project is a load balancer built with Golang. It is a program that keeps track of which server has the most traffic and which server is available. It then decides which server to send traffic to. The load balancer uses the round robin method to route traffic between servers. This method is the simplest way to route traffic between servers. The load balancer also uses a reverse proxy server to redirect requests and hide the location of the server from the client.

---

### Load Balancer
A load balancer is a device or software application that distributes network or application traffic across multiple servers. This distribution helps to improve the performance, reliability, and scalability of web applications and services.
    ![Load Balancer img](/images/loadbalancer.jpg)

### Reverse proxy
A reverse proxy is a server that sits between clients and origin servers. It intercepts client requests, processes them, and forwards them to the appropriate origin server.
    ![Reverse Proxy](/images/reverse-proxy.jpg)

### Struct and Functions
1. Load Balancer Struct:

- Servers: An array of Server structs, containing information about the backend servers.
- Round Robin Count: An integer tracking the current server index for round-robin load balancing.
- Function: A function to create a new Load Balancer instance.

    ```go
    type loadBalancer struct {
        port            string
        roundRobinCount int
        servers         []Server
    }
    ```

2. Server Struct:

- Address: The address of the backend server.
- Proxy: A function to create a new proxy connection to the backend server.
- Function: A function to create a new Server instance.

    ```go
    type Server interface {
        Address() string
        IsAlive() bool
        Serve(rw http.ResponseWriter, r *http.Request)
    }
    ```

3. Create a new Load Balancer:
- This function initializes a new Load Balancer struct, setting up the initial state and potentially performing necessary configurations.

    ```go
    func newLoadBalancer(port string, servers []Server) *loadBalancer {
        return &loadBalancer{
            port:            port,
            roundRobinCount: 0,
            servers:         servers,
        }
    }
    ```

4. Create a new Server:
- This function initializes a new Server struct, specifying its address and potentially setting up the proxy connection.

    ```go
    func newSimpleServer(addr string) *simpleServer {
        serverUrl, err := url.Parse(addr)
        handleErr(err)

        return &simpleServer{
            addrs: addr,
            proxy: httputil.NewSingleHostReverseProxy(serverUrl),
        }
    }
    ```

    ![Struct](/images/struct_function.jpg)

### Server Interface
This diagram illustrates the `Server` interface, which defines the core functionalities required by a server in your load balancing system.

1. **Address()**: This method returns the address of the server, which can be a hostname, IP address, or a combination of both.
2. **isAlive()**: This method checks the health status of the server. It returns a boolean value indicating whether the server is currently operational and can handle incoming requests.
3. **Serve()**: This method handles incoming requests to the server. It processes the requests, performs the necessary actions, and sends appropriate responses back to the client.

    ```go
    func (s *simpleServer) Address() string { return s.addrs }

    func (s *simpleServer) IsAlive() bool { return true }

    func (s *simpleServer) Serve(rw http.ResponseWriter, r *http.Request) {
        s.proxy.ServeHTTP(rw, r)
    }
    ```

    ![Server Interface](/images/Methods.jpg)

### Load Balancer Workflow
1. Main Function:

- Initializes the load balancer and sets up necessary configurations.
- Creates a pool of available servers.
- Starts the main loop to handle incoming requests.

2. Server Proxy Function:

- Receives an incoming request.
- Calls the GetNextAvailableServer function to obtain the next available server.
- Forwards the request to the selected server.
- Handles the response from the server and sends it back to the client.

3. GetNextAvailableServer Function:

- Iterates through the pool of servers using a round-robin algorithm.
- Checks the health of each server using a health check mechanism (e.g., sending a ping or HTTP request).
- Returns the first healthy server found.
- If no healthy servers are available, it can either wait for a server to recover or return an error to the client.

    ![workflow img](/images/mainFlow.jpg)

--- 


## How it Works
1. **Backend Server Registration:** The load balancer keeps track of available backend servers.
2. **Health Checks:** Periodically checks the health of each backend server.
3. **Load Balancing:** Uses the round robin algorithm to distribute incoming traffic evenly among healthy servers.
4. **Request Forwarding:** Forwards incoming requests to the selected backend server.
5. **Response Handling:** Receives responses from the backend server and forwards them to the client.


---

## Getting Started
1. Clone the Repository:

    ```bash
    git clone https://github.com/dev-dhanushkumar/Golang-Projects.git
    ```

2. Change the Directory

    ```bash
    cd go-loadbalancer/src
    ```

3. Run the Load Balancer
    ```bash
    go run main.go
    ```

--- 

## Features

- **Load balancing**: Distributes traffic across multiple servers to improve performance and reliability.
- **Round robin algorithm**: A simple method for distributing traffic evenly among servers.
- **Reverse proxy**: Hides the location of the server from the client and provides additional security and performance benefits.
- **Golang**: A modern programming language that is well-suited for building network applications.

---
