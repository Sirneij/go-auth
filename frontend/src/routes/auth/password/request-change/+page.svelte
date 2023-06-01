<script>
	import { applyAction, enhance } from '$app/forms';
	import { receive, send } from '$lib/utils/helpers';
	import { scale } from 'svelte/transition';

	/** @type {import('./$types').ActionData} */
	export let form;

	/** @type {import('./$types').SubmitFunction} */
	const handleRequestChange = async () => {
		return async ({ result }) => {
			await applyAction(result);
		};
	};
</script>

<div class="container">
	<form class="content" method="POST" use:enhance={handleRequestChange}>
		<h1 class="step-title">Request Password Change</h1>
		{#if form?.errors}
			{#each form?.errors as error (error.id)}
				<h4
					class="step-subtitle warning"
					in:receive={{ key: error.id }}
					out:send={{ key: error.id }}
				>
					{error.error}
				</h4>
			{/each}
		{/if}

		<div class="input-box">
			<span class="label">Email:</span>
			<input
				class="input"
				type="email"
				name="email"
				id="email"
				placeholder="Verified e-mail address"
				required
			/>
		</div>
		{#if form?.fieldsError && form?.fieldsError.email}
			<p class="warning" transition:scale|local={{ start: 0.7 }}>
				{form?.fieldsError.email}
			</p>
		{/if}

		<button class="button-colorful">Request change</button>
	</form>
</div>
