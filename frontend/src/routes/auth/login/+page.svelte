<script>
	import { applyAction, enhance } from '$app/forms';
	import { page } from '$app/stores';
	import { receive, send } from '$lib/utils/helpers';

	/** @type {import('./$types').ActionData} */
	export let form;

	/** @type {import('./$types').SubmitFunction} */
	const handleLogin = async () => {
		return async ({ result }) => {
			await applyAction(result);
		};
	};

	let message = '';
	if ($page.url.search) {
		message = $page.url.search.split('=')[1].replaceAll('%20', ' ');
	}
</script>

<div class="container">
	<form class="content" method="POST" action="?/login" use:enhance={handleLogin}>
		<h1 class="step-title">Login User</h1>
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

		{#if message}
			<h4 class="step-subtitle">{message}</h4>
		{/if}

		<input type="hidden" name="next" value={$page.url.searchParams.get('next')} />
		<div class="input-box">
			<span class="label">Email:</span>
			<input class="input" type="email" name="email" placeholder="Email address" />
		</div>
		<div class="input-box">
			<span class="label">Password:</span>
			<input class="input" type="password" name="password" placeholder="Password" />
			<a href="/#" style="margin-left: 1rem;">Forgot password?</a>
		</div>
		<div class="btn-container">
			<button class="button-colorful">Login</button>
			<p>Have no account? <a href="/auth/register">Register here</a>.</p>
		</div>
	</form>
</div>
