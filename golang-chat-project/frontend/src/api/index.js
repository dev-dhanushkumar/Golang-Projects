var socket  = new WebSocket['ws://localhost:9000/ws'];

let connect = (cb) => {
    console.log("Connecting...");

    socket.onopen  = () => {
        console.log("Successfully connected");
    }

    socket.onmessage = (msg) => {
        console.log("Message from websocket:", msg);
    }

    socket.onclose = (event) => {
        console.log("Socket closed connection:", event);
    }

    socket.onerror = (error) => {
        console.log("Socket error:", error);
    }
};

let sendMsg =  (msg) => {
    console.log("Sending Message: ", msg);
    socket.send(msg);
}

export { connect, sendMsg };