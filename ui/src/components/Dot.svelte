<script>
    export let state = false;
</script>

<main class:state></main>

<style>
    main {
        height: 12px;
        width: 12px;
        border-radius: 50%;
        display: inline-block;
        background-color: red;
    }

    .state {
        background-color: limegreen;
    }
</style>
