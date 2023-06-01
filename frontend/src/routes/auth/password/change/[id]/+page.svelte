<script>
	import { applyAction, enhance } from '$app/forms';
	import { page } from '$app/stores';
	import { receive, send } from '$lib/utils/helpers';
	import { scale } from 'svelte/transition';

	/** @type {import('./$types').ActionData} */
	export let form;

	/** @type {import('./$types').SubmitFunction} */
	const handleChange = async () => {
		return async ({ result }) => {
			await applyAction(result);
		};
	};
</script>

<div class="container">
	<form class="content" method="POST" use:enhance={handleChange}>
		<h1 class="step-title title">Change Password</h1>
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
				pattern="\d*"
				maxlength="6"
				minlength="6"
			/>
		</div>
		<div class="input-box">
			<span class="label">Password:</span>
			<input class="input" type="password" name="password" placeholder="Password" />
		</div>
		{#if form?.fieldsError && form?.fieldsError.password}
			<p class="warning" transition:scale|local={{ start: 0.7 }}>
				{form?.fieldsError.password}
			</p>
		{/if}
		<div class="input-box">
			<span class="label">Confirm password:</span>
			<input class="input" type="password" name="confirm_password" placeholder="Password" />
		</div>
		{#if form?.fieldsError && form?.fieldsError.confirmPassword}
			<p class="warning" transition:scale|local={{ start: 0.7 }}>
				{form?.fieldsError.confirmPassword}
			</p>
		{/if}

		<button class="button-colorful">Change password</button>
	</form>
</div>
