    let socket = null;
    let o = document.getElementById("output");

    window.addEventListener("beforeunload", function (e) {
        disconnectClient();
    });


    document.addEventListener("DOMContentLoaded", function() {
        socket = new ReconnectingWebSocket("ws://127.0.0.1:8080/ws", null, { debug: true, reconnectInterval: 3000 });
        
        const OFFLINE = `<span class="badge bg-danger">not connected</span>`;
        const ONLINE = `<span class="badge bg-success">connected</span>`;
        const ERROR = `<span class="badge bg-warning">error</span>`;
        let statusDiv = document.getElementById("status");

        socket.onopen = () => {
            console.log("connected");
            statusDiv.innerHTML = ONLINE;
        }

        socket.onclose = () => {
            console.log("connection closed");
            statusDiv.innerHTML = OFFLINE;
        }

        socket.onerror = error => {
            console.log("error" + JSON.stringify(error));
            statusDiv.innerHTML = ERROR;
        }

        socket.onmessage = msg => {

            let data = JSON.parse(msg.data);
            //console.log(data);
            //console.log("Action:", data.action)
            switch (data.action) {
                case "list_users":
                    let ul = document.getElementById("online_users");
                    ul.innerHTML = "";

                    if (data.connected_users.length > 0) {
                        data.connected_users.forEach(function(item) {
                            let li = document.createElement("li");
                            li.appendChild(document.createTextNode(item));
                            ul.appendChild(li);
                        });
                    }
                break;
                case "broadcast": 
                    console.log("new message:", data.message);
                    newD = document.createElement("div");
                    newD.innerHTML = data.message;
                    o.appendChild(newD);
                break;
            }
        }

        let userInput = document.getElementById("username");
        userInput.addEventListener("change", function() {
            let jsonData = {};
            jsonData["action"] = "username";
            jsonData["username"] = this.value;
            socket.send(JSON.stringify(jsonData));
        });
    });

    document.getElementById("message").addEventListener("keydown", function(e) {
        if (e.code === "Enter") {
            if (!socket) {
                console.log("not connected");
                return false;
            }
            e.preventDefault();
            e.stopPropagation();
            if (checkFields()) sendMessage();
        }
    });

    document.getElementById("sendBtn").addEventListener("click", function(e) {
        if (checkFields()) {sendMessage();}
    });

    function checkFields() {
        const checkIDs = ["username", "message"];
        let checked = true;
        checkIDs.forEach(function(id) {
            if (!checkFieldIsNotEmbyByID(id)) {
                checked = false;
            }
        });
        return checked;
    }

    function checkFieldIsNotEmbyByID(id, className = "alert-warning") {
        let field = document.getElementById(id);
        field.classList.remove(className);
        if (field.value === "") {
            field.classList.add(className);
            return false;
        }
        return true;
    }

    function sendMessage() {
        let jsonData = {};
        jsonData["action"] = "broadcast";
        jsonData["username"] = document.getElementById("username").value;
        let message = document.getElementById("message");
        jsonData["message"] = message.value;
        socket.send(JSON.stringify(jsonData));
        message.value = "";
    }

    function disconnectClient() {
        console.log("dc client");
        let jsonData = {};
        jsonData["action"] = "left";
        socket.send(JSON.stringify(jsonData));
    }