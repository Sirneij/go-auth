<script>
	/** @type {number | null} */
	export let width;
	/** @type {string | null} */
	export let message;
</script>

<div class="loading">
	<p class="simple-loader" style={width ? `width: ${width}px` : ''} />
	{#if message}
		<p>{message}</p>
	{/if}
</div>

<style>
	.loading {
		display: flex;
		align-items: center;
		justify-content: center;
	}
	.loading p {
		margin-left: 0.5rem;
	}
</style>
