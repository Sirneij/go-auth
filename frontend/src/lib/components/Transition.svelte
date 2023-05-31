<script>
	import { slide } from 'svelte/transition';
	/** @type {string} */
	export let key;

	/** @type {number} */
	export let duration = 300;
</script>

{#key key}
	<div in:slide={{ duration, delay: duration }} out:slide={{ duration }}>
		<slot />
	</div>
{/key}
