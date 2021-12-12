<script>
    import MatchStepper from "./components/MatchStepper.svelte";
    import {onMount} from "svelte";
    import FieldTeam from "./components/FieldTeam.svelte";

    // let wsServer = "ws://" + location.host + "/ws";
    let wsServer = "ws://localhost:8080/ws";

    let ws;
    let startTime;
    let hideFTATools = true;

    let latency;
    let wsConnected = false;
    let matchState = {};
    let banner = "Waiting for FMS connection";

    function wsConnect() {
        ws = new WebSocket(wsServer)

        ws.onopen = () => {
            wsConnected = true
            banner = "FMS connection established"
        }
        ws.onclose = () => {
            wsConnected = false
            banner = "Retrying FMS connection..."
            setTimeout(function () {
                wsConnect()
            }, 1000)
        }
        ws.onerror = () => {
            wsConnected = false
            banner = "FMS connection error"
            ws.close()
        }
        ws.onmessage = (event) => {
            latency = Date.now() - startTime;
            matchState = JSON.parse(event.data)
            if (matchState["state"] === "Idle") {
                banner = "Ready to start match"
            } else { // Match running
                banner = "Running: " + matchState["state"]
            }
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

    function testSounds() {
        ws.send(JSON.stringify({
            message: "test_sounds"
        }))
    }

    function estop(teamNumber, allianceStation) {
        if (confirm(`Confirm E-STOP ${teamNumber} (${allianceStation})?`)) {
            ws.send(JSON.stringify({
                message: "estop",
                arg: teamNumber
            }))
            alert("E-stopped " + teamNumber)
        } else {
            alert("E-STOP Cancelled")
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

    onMount(() => {
        wsConnect()
        setInterval(function () {
            startTime = Date.now();
            ws.send(JSON.stringify({
                message: "ping"
            }));
        }, 1000)
    })
</script>

<main>
    <div>
        <h1>BunnyFMS</h1>
        <p class="fms-dot">
            FMS <span
                class="dot"
                style="background-color: {wsConnected ? 'green' : 'red'}"></span>
        </p>
    </div>

    <h2>{banner}</h2>

    <MatchStepper
            autoTimer={matchState["auto_timer"]}
            endgameTimer={matchState["endgame_timer"]}
            state={matchState["state"]}
            teleopTimer={matchState["teleop_timer"]}
    />

    <div class="field">
        <div class="alliance">
            {#each ["R1", "R2", "R3"] as allianceStation}
                <FieldTeam {estop} disabled={!matchState['state'] || matchState['state'] === "Idle"} allianceStation={allianceStation} teamNumber="0000"/>
            {/each}
        </div>

        <div class="match-controls">
            {#if matchState && matchState['state']}
                {#if matchState['state'] === "Idle"}
                    <button on:click={() => startMatch()}>Start Match</button>
                {:else}
                    <button on:click={() => stopMatch()}>Stop Match</button>
                {/if}
            {/if}
        </div>

        <div class="alliance text-align-right">
            {#each ["B1", "B2", "B3"] as allianceStation}
                <FieldTeam {estop} disabled={!matchState['state'] || matchState['state'] === "Idle"} allianceStation={allianceStation} teamNumber="0000"/>
            {/each}
        </div>
    </div>

    <p on:click={() => {hideFTATools = !hideFTATools}}>FTA Tools ▼</p>
    {#if !hideFTATools}
        <div class="hidden" id="fta-tools">
            <p>WS latency: {latency} ms</p>
            <button on:click={() => wsConnect()}>WS Refresh</button>
            <button on:click={() => dsReconnect()}>Force DS Reconnect</button>
            <button on:click={() => testSounds()}>Test game sounds</button>
        </div>
    {/if}
</main>

<style>
    .dot {
        height: 12px;
        width: 12px;
        border-radius: 50%;
        display: inline-block;
    }

    h2 {
        display: flex;
        justify-content: center;
        margin-top: 0;
    }

    .fms-dot {
        margin-bottom: 0;
        margin-top: 0;
    }

    .field {
        display: flex;
        justify-content: space-between;
        border-radius: 5px;
        padding-left: 25px;
        padding-right: 25px;
        padding-bottom: 25px;
        margin-bottom: 15px;
    }

    @media (prefers-color-scheme: light) {
        .field {
            border: 2px solid black;
        }
    }

    @media (prefers-color-scheme: dark) {
        .field {
            border: 2px solid white;
        }
    }

    .text-align-right {
        text-align: right;
    }

    .match-controls {
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
    }
</style>
