<script>
	import { applyAction, enhance } from '$app/forms';
	import { page } from '$app/stores';

	import { receive, send } from '$lib/utils/helpers';

	/** @type {import('./$types').ActionData} */
	export let form;

	/** @type {import('./$types').SubmitFunction} */
	const handleActivate = async () => {
		return async ({ result }) => {
			await applyAction(result);
		};
	};

	/** @type {string|undefined} */
	let token = '';

	$: if (token) {
		token = token.split(' ').join('');
		let finalVal = token.match(/.{1,3}/g)?.join(' ');
		token = finalVal;
	}
</script>

<div class="container">
	<form class="content" method="POST" use:enhance={handleActivate}>
		<h1 class="step-title">Activate your account</h1>
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

		<input type="hidden" name="user_id" value={$page.params.id} />
		<div class="input-box">
			<span class="label">Token:</span>
			<input
				type="text"
				class="input"
				name="token"
				placeholder="XXX XXX"
				inputmode="numeric"
				bind:value={token}
				maxlength="7"
				minlength="7"
			/>
		</div>

		<button class="button-dark">Activate</button>
	</form>
</div>
