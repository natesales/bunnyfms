<script>
    import Dot from "./Dot.svelte";

    export let estop, updateAlliances, editTeamNumbers;
    export let matchIdle = false;
    export let allianceStation, teamNumber;

    let isBlueAlliance = allianceStation.startsWith('B');
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
    <!--    <p>-->
    <!--        DS:-->
    <!--        <Dot state={false}/>-->
    <!--    </p>-->
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
</style>
