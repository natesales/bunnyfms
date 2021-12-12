<script>
    import FieldTeam from "./FieldTeam.svelte";
    import {ws} from "../api";

    export let matchState;

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
</script>

<main>
    <div class="alliance">
        {#each ["R1", "R2", "R3"] as allianceStation}
            <FieldTeam disabled={!matchState['state'] || matchState['state'] === "Idle"} allianceStation={allianceStation} teamNumber="0000"/>
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
            <FieldTeam disabled={!matchState['state'] || matchState['state'] === "Idle"} allianceStation={allianceStation} teamNumber="0000"/>
        {/each}
    </div>
</main>

<style>
    main {
        display: flex;
        justify-content: space-between;
        border-radius: 5px;
        padding-left: 25px;
        padding-right: 25px;
        padding-bottom: 25px;
        margin-bottom: 15px;
    }

    @media (prefers-color-scheme: light) {
        main {
            border: 2px solid black;
        }
    }

    @media (prefers-color-scheme: dark) {
        main {
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
