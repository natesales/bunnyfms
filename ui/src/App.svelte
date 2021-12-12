<script>
    import {onMount} from "svelte";
    import FieldTeam from "./components/FieldTeam.svelte";
    import Dot from "./components/Dot.svelte";

    // let wsServer = "ws://" + location.host + "/ws";
    let wsServer = "ws://localhost:8080/ws";

    let ws;
    let startTime;
    let hideFTATools = true;

    let latency;
    let wsConnected = false;
    let matchState = {};
    let banner = "Waiting for FMS connection";
    let allianceMap = {};
    let editingTeamNumbers = false;
    let editingMatchName = false;
    let matchName;
    let readyToStart = false;

    // https://stackoverflow.com/questions/5072136/javascript-filter-for-objects/37616104
    Object.filter = (obj, predicate) =>
        Object.keys(obj)
            .filter(key => predicate(obj[key]))
            .reduce((res, key) => (res[key] = obj[key], res), {});

    function editTeamNumbers() {
        editingTeamNumbers = true
    }

    function wsConnect() {
        ws = new WebSocket(wsServer)

        ws.onopen = () => {
            wsConnected = true
            banner = "FMS connection established"
        }
        ws.onclose = () => {
            wsConnected = false
            banner = "Lost FMS connection"
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
                readyToStart = false;

                // Check if each alliance has at least one team and all configured teams' drive stations have connected
                let hasRed = false;
                let hasBlue = false;
                let waitingFor = [];
                for (let position in matchState["alliances"]) {
                    let teamNumber = matchState["alliances"][position]
                    if (teamNumber > 0) {
                        if (position.startsWith("R")) {
                            hasRed = true
                        } else if (position.startsWith("B")) {
                            hasBlue = true
                        }
                        if (!matchState["ds"] || !matchState["ds"][position]) {
                            waitingFor.push(teamNumber)
                        }
                    }
                }

                if (!(hasRed && hasBlue)) {
                    banner = "Ready to configure match"
                } else if (waitingFor.length !== 0) {
                    banner = "Waiting for " + waitingFor.length + " team"
                    if (waitingFor.length > 1) {
                        banner += "s"
                    }
                } else if (matchName === "") {
                    banner = "Please set a match name"
                } else {
                    banner = "Ready to start match"
                    readyToStart = true
                }

                if (!editingTeamNumbers) {
                    allianceMap = Object.filter(matchState["alliances"], x => (x && x !== 0))
                }
                if (!editingMatchName) {
                    matchName = matchState["name"]
                }
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
        if (confirm("Are you sure you want to test game sounds?")) {
            ws.send(JSON.stringify({
                message: "test_sounds"
            }))
            alert("Playing game sounds")
        } else {
            alert("Game sound test cancelled")
        }
    }

    function resetAlliances() {
        if (confirm("Are you sure you want to reset alliances?")) {
            ws.send(JSON.stringify({
                message: "reset_alliances"
            }))
            alert("Alliances reset")
        } else {
            alert("Alliance reset cancelled")
        }
    }

    function estop(teamNumber, allianceStation) {
        if (confirm(`Confirm E-STOP ${teamNumber} (${allianceStation})?`)) {
            ws.send(JSON.stringify({
                message: "estop",
                alliance_station: allianceStation
            }))
            alert("E-stopped " + teamNumber)
        } else {
            alert("E-STOP Cancelled")
        }
    }

    function startMatch() {
        ws.send(JSON.stringify({
            message: "start"
        }))
    }

    function stopMatch() {
        ws.send(JSON.stringify({
            message: "stop"
        }))
    }

    function updateAlliances() {
        allianceMap = Object.filter(allianceMap, x => (x && x !== 0))

        ws.send(JSON.stringify({
            message: "update_alliances",
            alliances: allianceMap
        }))
        editingTeamNumbers = false
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
    <div class="space-between">
        <h2>BunnyFMS</h2>
        {#if matchState["event_name"]}
            <h2>{matchState["event_name"]}</h2>
        {/if}
    </div>
    <div class="field">
        <div class="alliance">
            <FieldTeam allianceStation="R1" bind:matchState={matchState} bind:teamNumber={allianceMap["R1"]} {editTeamNumbers} {estop} {updateAlliances}/>
            <FieldTeam allianceStation="R2" bind:matchState={matchState} bind:teamNumber={allianceMap["R2"]} {editTeamNumbers} {estop} {updateAlliances}/>
            <FieldTeam allianceStation="R3" bind:matchState={matchState} bind:teamNumber={allianceMap["R3"]} {editTeamNumbers} {estop} {updateAlliances}/>
        </div>

        <div class="match-center">
            <h2 style="margin-bottom: 10px">{banner}</h2>

            {#if matchState['state']}
                <input
                        placeholder="Match name"
                        style="text-align: center"
                        type="text"
                        disabled={matchState["state"] !== "Idle"}
                        bind:value={matchName}
                        on:focus={() => editingMatchName=true}
                        on:blur={() => {
                            ws.send(JSON.stringify({
                                message: "match_name",
                                name: matchName
                            }))
                            editingMatchName=false
                        }}
                >
                <h2 style="margin-bottom: 0">{matchState["current_timer"]}</h2>
                <div class="match-timers">
                    <p>Auto: {matchState["auto_timer"]}</p>
                    <p>Teleop: {matchState["teleop_timer"]}</p>
                    <p>Endgame: {matchState["endgame_timer"]}</p>
                </div>

                {#if matchState['state'] === "Idle"}
                    <button disabled={!readyToStart} on:click={() => startMatch()}>Start Match</button>
                {:else}
                    <button on:click={() => stopMatch()}>Stop Match</button>
                {/if}
            {/if}
        </div>

        <div class="alliance text-align-right">
            <FieldTeam allianceStation="B1" bind:matchState={matchState} bind:teamNumber={allianceMap["B1"]} {editTeamNumbers} {estop} {updateAlliances}/>
            <FieldTeam allianceStation="B2" bind:matchState={matchState} bind:teamNumber={allianceMap["B2"]} {editTeamNumbers} {estop} {updateAlliances}/>
            <FieldTeam allianceStation="B3" bind:matchState={matchState} bind:teamNumber={allianceMap["B3"]} {editTeamNumbers} {estop} {updateAlliances}/>
        </div>
    </div>

    <div class="space-between">
        <p on:click={() => {hideFTATools = !hideFTATools}}>FTA Tools ▼</p>
        <p class="fms-dot">
            FMS:
            <Dot state={wsConnected}/>
        </p>
    </div>
    {#if !hideFTATools}
        <div class="hidden" id="fta-tools">
            <p>WS latency: {latency} ms</p>
            <button on:click={() => wsConnect()}>WS Refresh</button>
            <button on:click={() => dsReconnect()}>Force DS Reconnect</button>
            <button on:click={() => testSounds()}>Test game sounds</button>
            <button on:click={() => resetAlliances()}>Reset alliances</button>
        </div>
    {/if}
</main>

<style>
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

    .match-center {
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
    }

    .match-timers {
        margin-top: 5px;
        margin-bottom: 14px;
        display: flex;
        align-items: center;
        flex-direction: column;
    }

    .match-timers p {
        margin-top: 2px;
        margin-bottom: 2px;
    }

    .space-between {
        display: flex;
        flex-direction: row;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 5px;
    }

    .space-between p, h2 {
        margin: 0;
        padding: 0;
    }
</style>
