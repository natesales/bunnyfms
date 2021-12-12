<script>
    export let disabled = true;
    export let allianceStation, teamNumber;

    function estop() {
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
</script>

<main>
    <h2>
        <span style="color: {allianceStation.startsWith('R') ?'red' : 'blue'}">
            {allianceStation}
        </span>
        {teamNumber}
    </h2>
    <button {disabled} on:click={() => {estop(allianceStation)}}>E-STOP</button>
</main>

<style>
    button {
        font-weight: bold;
        background-color: red;
    }
</style>
