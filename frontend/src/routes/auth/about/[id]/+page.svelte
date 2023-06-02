<script>
	import { applyAction, enhance } from '$app/forms';
	import { page } from '$app/stores';
	import ImageInput from '$lib/components/ImageInput.svelte';
	import Modal from '$lib/components/Modal.svelte';
	import Avatar from '$lib/img/teamavatar.png';
	import { receive, send } from '$lib/utils/helpers';

	let showModal = false;
	const open = () => (showModal = true);
	const close = () => (showModal = false);

	/** @type {import('./$types').ActionData} */
	export let form;

	/** @type {import('./$types').SubmitFunction} */
	const handleUpdate = async () => {
		return async ({ result }) => {
			if (result.type === 'success' || result.type === 'redirect') {
				close();
			}
			await applyAction(result);
		};
	};
</script>

<div class="hero-container">
	<div class="hero-logo">
		<img
			src={$page.data.user.thumbnail ? $page.data.user.thumbnail : Avatar}
			alt={`${$page.data.user.first_name} ${$page.data.user.last_name}`}
		/>
	</div>
	<h3 class="hero-subtitle subtitle">
		Name (First and Last): {`${$page.data.user.first_name} ${$page.data.user.last_name}`}
	</h3>
	{#if $page.data.user.profile.phone_number}
		<h3 class="hero-subtitle">
			Phone: {$page.data.user.profile.phone_number}
		</h3>
	{/if}

	{#if $page.data.user.profile.github_link}
		<h3 class="hero-subtitle">
			GitHub: {$page.data.user.profile.github_link}
		</h3>
	{/if}

	{#if $page.data.user.profile.birth_date}
		<h3 class="hero-subtitle">
			Date of birth: {$page.data.user.profile.birth_date}
		</h3>
	{/if}
	<div class="hero-buttons-container">
		<button class="button-dark" on:click={open}>Edit profile</button>
	</div>
</div>

{#if showModal}
	<Modal on:close={close}>
		<form class="content" method="POST" enctype="multipart/form-data" use:enhance={handleUpdate}>
			<h1 class="step-title" style="text-align: center;">Login User</h1>
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

			<ImageInput avatar={$page.data.user.thumbnail} />

			<div class="input-box">
				<span class="label">First name:</span>
				<input
					class="input"
					type="text"
					name="first_name"
					value={$page.data.user.first_name}
					placeholder="Your first name..."
				/>
			</div>
			<div class="input-box">
				<span class="label">Last name:</span>
				<input
					class="input"
					type="text"
					name="last_name"
					value={$page.data.user.last_name}
					placeholder="Your last name..."
				/>
			</div>
			<div class="input-box">
				<span class="label">Phone number:</span>
				<input
					class="input"
					type="tel"
					name="phone_number"
					value={$page.data.user.profile.phone_number ? $page.data.user.profile.phone_number : ''}
					placeholder="Your phone number e.g +2348135703593..."
				/>
			</div>
			<div class="input-box">
				<span class="label">Birth date:</span>
				<input
					class="input"
					type="date"
					name="birth_date"
					value={$page.data.user.profile.birth_date ? $page.data.user.profile.birth_date : ''}
					placeholder="Your date of birth..."
				/>
			</div>
			<div class="input-box">
				<span class="label">GitHub Link:</span>
				<input
					class="input"
					type="url"
					name="github_link"
					value={$page.data.user.profile.github_link ? $page.data.user.profile.github_link : ''}
					placeholder="Your github link e.g https://github.com/Sirneij/..."
				/>
			</div>
			<button type="submit" class="button-colorful">Update</button>
		</form>
	</Modal>
{/if}

<style>
	.hero-container .hero-subtitle:not(:last-of-type) {
		margin: 0 0 0 0;
	}
</style>
