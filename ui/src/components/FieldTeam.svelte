<script>
    export let estop;
    export let disabled = true;
    export let allianceStation, teamNumber;
</script>

<main>
    <h2>
        <span style="color: {allianceStation.startsWith('R') ?'red' : 'blue'}">
            {allianceStation}
        </span>
        {teamNumber}
    </h2>
    <button {disabled} on:click={() => {estop(teamNumber, allianceStation)}}>E-STOP</button>
</main>

<style>
    button {
        font-weight: bold;
        background-color: red;
    }
</style>
