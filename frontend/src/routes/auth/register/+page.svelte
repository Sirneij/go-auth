<script>
	import { applyAction, enhance } from '$app/forms';
	import { receive, send } from '$lib/utils/helpers';
	import { scale } from 'svelte/transition';

	/** @type {import('./$types').ActionData} */
	export let form;

	/** @type {import('./$types').SubmitFunction} */
	const handleRegister = async () => {
		return async ({ result }) => {
			await applyAction(result);
		};
	};
</script>

<div class="container">
	<form class="content" action="?/register" method="POST" use:enhance={handleRegister}>
		<h1 class="step-title title">Register</h1>
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
			<input class="input" type="email" name="email" placeholder="Email address" />
		</div>
		{#if form?.fieldsError && form?.fieldsError.email}
			<p class="warning" transition:scale|local={{ start: 0.7 }}>
				{form?.fieldsError.email}
			</p>
		{/if}
		<div class="input-box">
			<span class="label">First name:</span>
			<input class="input" type="text" name="first_name" placeholder="First name" />
		</div>
		<div class="input-box">
			<span class="label">Last name:</span>
			<input class="input" type="text" name="last_name" placeholder="Last name" />
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

		<div class="btn-container">
			<button class="button-colorful">Register</button>
			<p>Already registered? <a href="/auth/login">Login here</a>.</p>
		</div>
	</form>
</div>
