<script>
    import Field from "./components/Field.svelte";
    import MatchStepper from "./components/MatchStepper.svelte";
    import {onMount} from "svelte";
    import {ws} from "./api";

    // let wsServer = "ws://" + location.host + "/ws";
    let wsServer = "ws://localhost:8080/ws";

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

    <Field {matchState}/>

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
</style>
