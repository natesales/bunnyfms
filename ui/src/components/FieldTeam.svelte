<script>
    import Dot from "./Dot.svelte";

    export let matchState;
    export let estop, updateAlliances, editTeamNumbers;
    export let allianceStation, teamNumber;

    let isBlueAlliance = false;
    let matchIdle = true;
    $:{
        isBlueAlliance = allianceStation.startsWith('B');
        matchIdle = !matchState['state'] || matchState['state'] === "Idle";
    }
</script>

<main>
    <input
            bind:value={teamNumber}
            class:align-right={isBlueAlliance}
            class:blue={isBlueAlliance}
            disabled={!matchIdle}
            on:blur={updateAlliances}
            on:focus={editTeamNumbers}
            type="number"
    >

    <p>
        {#if matchState["ds"] && matchState["ds"][allianceStation]}
            DS:
            <Dot state={matchState["ds"][allianceStation]["ds_link"]}/>
            (last packet {matchState["ds"][allianceStation]["last_packet"]})
            <br>
            Robot:
            <Dot state={matchState["ds"][allianceStation]["robot_link"]}/>
            (last link {matchState["ds"][allianceStation]["last_robot_link"]})
            <br>
            Robot:
            <Dot state={matchState["ds"][allianceStation]["radio_link"]}/>
            ({matchState["ds"][allianceStation]["battery_voltage"]}v) <span style="color: red; font-weight: bold">{matchState["ds"][allianceStation]["estop"] ? "E-STOPPED" : ""}</span>
            <br>
        {/if}
    </p>
    <button class:align-right={isBlueAlliance} disabled={matchIdle} on:click={() => {estop(teamNumber, allianceStation)}}>E-STOP</button>
</main>

<style>
    main {
        margin-top: 20px;
    }

    button {
        font-weight: bold;
        background-color: #ee1b1b;
    }

    input {
        margin-top: 10px;
        margin-bottom: 10px;
        width: 4ch;
        border: 2px solid red;
    }

    .blue {
        border: 2px solid blue;
    }

    /* Hide number arrows */
    input::-webkit-outer-spin-button,
    input::-webkit-inner-spin-button {
        -webkit-appearance: none;
        margin: 0;
    }

    input[type=number] {
        -moz-appearance: textfield;
    }

    .align-right {
        margin-left: auto;
        margin-right: 0;
    }
</style>
