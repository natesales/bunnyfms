<script>
    export let state;
    export let autoTimer = "-";
    export let teleopTimer = "-";
    export let endgameTimer = "-";
</script>

<main>
    <ul class="stepper">
        <li class:active={state === "Idle"}>Idle</li>
        <li class:active={state === "Auto" || state === "Teleop" || state === "Endgame"}>Auto: {autoTimer}</li>
        <li class:active={state === "Teleop" || state === "Endgame"}>Teleop: {teleopTimer}</li>
        <li class:active={state === "Endgame"}>Endgame: {endgameTimer}</li>
    </ul>
</main>

<style>
    .stepper li {
        list-style-type: none;
        float: left;
        width: 25%;
        position: relative;
        text-align: center;
        margin-bottom: 30px;
    }

    .stepper li:after {
        content: "";
        position: absolute;
        width: 100%;
        height: 4px;
        background-color: #ddd;
        top: 25px;
        left: -50%;
    }

    .stepper li {
        color: white;
    }

    @media (prefers-color-scheme: light) {
        .stepper li {
            color: black;
        }
    }

    @media (prefers-color-scheme: dark) {
        .stepper li {
            color: white;
        }
    }

    .stepper li:first-child:after {
        content: none;
    }

    .stepper li.active:before {
        border-color: grey;
        background-color: grey;
    }

    .stepper .active:after {
        background-color: grey;
    }
</style>
