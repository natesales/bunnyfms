let ws;
let startTime;

function setBanner(text) {
    document.getElementById("banner").innerText = text
}

function loadTeams() {
    r1 = document.getElementById("r1")
    r2 = document.getElementById("r2")
    r3 = document.getElementById("r3")
    b1 = document.getElementById("b1")
    b2 = document.getElementById("b2")
    b3 = document.getElementById("b3")

    r1.innerText = "0000"
    r2.innerText = "0000"
    r3.innerText = "0000"
    b1.innerText = "0000"
    b2.innerText = "0000"
    b3.innerText = "0000"
}

function setTimers(state, autoTime, teleopTime, endgameTime) {
    document.getElementById("timer-auto").innerText = ""
    document.getElementById("timer-teleop").innerText = ""
    document.getElementById("timer-endgame").innerText = ""

    if (state === "Auto") {
        document.getElementById("timer-auto").innerText = ": " + autoTime
    } else if (state === "Teleop") {
        document.getElementById("timer-teleop").innerText = ": " + teleopTime
    } else if (state === "Endgame") {
        document.getElementById("timer-teleop").innerText = ": " + teleopTime
        document.getElementById("timer-endgame").innerText = ": " + endgameTime
    }
}

function updateLatency(latency) {
    document.getElementById("latency").innerText = latency + " ms";
}

function wsError(message) {
    document.getElementById("connection-indicator").classList.remove("bg-green")
    document.getElementById("connection-indicator").classList.add("bg-red")
    updateLatency("-")
    setBanner("FMS backend connection lost" + (message ? ": " + message : ""))
}

function updateStepper(state) {
    if (state === "Idle") {
        document.getElementById("stepper-idle").className = "active"
        document.getElementById("stepper-auto").className = ""
        document.getElementById("stepper-teleop").className = ""
        document.getElementById("stepper-endgame").className = ""
    } else if (state === "Auto") {
        document.getElementById("stepper-idle").className = "active"
        document.getElementById("stepper-auto").className = "active"
        document.getElementById("stepper-teleop").className = ""
        document.getElementById("stepper-endgame").className = ""
    } else if (state === "Teleop") {
        document.getElementById("stepper-idle").className = "active"
        document.getElementById("stepper-auto").className = "active"
        document.getElementById("stepper-teleop").className = "active"
        document.getElementById("stepper-endgame").className = ""
    } else if (state === "Endgame") {
        document.getElementById("stepper-idle").className = "active"
        document.getElementById("stepper-auto").className = "active"
        document.getElementById("stepper-teleop").className = "active"
        document.getElementById("stepper-endgame").className = "active"
    }
}

function wsConnect() {
    ws = new WebSocket("ws://" + location.host + "/ws")

    ws.onopen = () => {
        document.getElementById("connection-indicator").classList.remove("bg-red")
        document.getElementById("connection-indicator").classList.add("bg-green")
        setBanner("FMS connection established")
    }
    ws.onclose = (event) => {
        wsError(event.reason)
        setTimeout(function () {
            wsConnect()
        }, 1000)
    }
    ws.onerror = (event) => {
        wsError(event.message)
        ws.close()
    }
    ws.onmessage = (event) => {
        let latency = Date.now() - startTime;
        let body = JSON.parse(event.data)
        if (body["state"] === "Idle") {
            setBanner("Ready to start match")
            document.getElementById("start-match").classList.remove("hidden")
            document.getElementById("stop-match").classList.add("hidden")
        } else { // Match running
            document.getElementById("start-match").classList.add("hidden")
            document.getElementById("stop-match").classList.remove("hidden")
            setBanner("Running: " + body["state"])
        }
        setTimers(body["state"], body["auto_timer"], body["teleop_timer"], body["endgame_timer"])
        updateStepper(body["state"])
        updateLatency(latency)
    }
}

function startMatch() {
    ws.send(JSON.stringify({
        message: "start"
    }));
}

function stopMatch() {
    ws.send(JSON.stringify({
        message: "stop"
    }));
}


function estop(team) {
    if (confirm(`Confirm ESTOP ${team} (${document.getElementById("r1").innerText})?`)) {
        ws.send(JSON.stringify({
            message: "estop",
            arg: team
        }))
        alert("E-stopped " + team)
    } else {
        alert("ESTOP Cancelled")
    }
}

function dsReconnect() {
    if (confirm("Are you sure you want to force a DS reconnect?")) {
        ws.send(JSON.stringify({
            message: "ds_reconnect"
        }));
        alert("Sent DS reconnect")
    } else {
        alert("DS reconnect cancelled")
    }
}

window.addEventListener("DOMContentLoaded", (event) => {
    loadTeams()
    wsConnect()

    setInterval(function () {
        startTime = Date.now();
        ws.send(JSON.stringify({
            message: "ping"
        }));
    }, 1000)
})
