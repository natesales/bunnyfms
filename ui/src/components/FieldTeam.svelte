<script>
    import Dot from "./Dot.svelte";

    export let estop, updateAlliances, editTeamNumbers;
    export let matchIdle = false;
    export let allianceStation, teamNumber;

    let rightAlign = allianceStation.startsWith('B');
</script>

<main>
    <div class="row" class:align-right={rightAlign}>
        {#if rightAlign}
            <h2 class:align-right={rightAlign} style="color: {allianceStation.startsWith('R') ?'red' : 'blue'}">
                {allianceStation}
            </h2>
        {/if}
        <input
                bind:value={teamNumber}
                class:align-right={rightAlign}
                disabled={!matchIdle}
                on:blur={updateAlliances}
                on:focus={editTeamNumbers}
                type="number"
        >
        {#if !rightAlign}
            <h2 style="color: {allianceStation.startsWith('R') ?'red' : 'blue'}">
                {allianceStation}
            </h2>
        {/if}
    </div>
    <p>
        DS:
        <Dot state={false}/>
    </p>
    <button class:align-right={allianceStation.startsWith('B')} disabled={matchIdle} on:click={() => {estop(teamNumber, allianceStation)}}>E-STOP</button>
</main>

<style>
    main {
        margin-top: 20px;
    }

    button {
        font-weight: bold;
        background-color: red;
    }

    input {
        margin-top: 10px;
        margin-bottom: 10px;
        width: 5ch;
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

    .align-right input {
        margin-right: 0;
        padding-right: 0;
        margin-left: 5px;
    }

    p {
        margin-top: 0;
        margin-bottom: 10px;
        padding: 0;
    }

    .row {
        display: flex;
        flex-direction: row;
        align-items: center;
    }

    .row h2 {
        margin-top: 12px;
        margin-bottom: 12px;
    }
</style>
