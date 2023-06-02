<script>
	import { applyAction, enhance } from '$app/forms';
	import { page } from '$app/stores';
	import ImageInput from '$lib/components/ImageInput.svelte';
	import Modal from '$lib/components/Modal.svelte';
	import SmallLoader from '$lib/components/SmallLoader.svelte';
	import Avatar from '$lib/img/teamavatar.png';
	import { receive, send } from '$lib/utils/helpers';

	$: ({ user } = $page.data);

	let showModal = false,
		isUploading = false,
		isUpdating = false;
	const open = () => (showModal = true);
	const close = () => (showModal = false);

	/** @type {import('./$types').ActionData} */
	export let form;

	/** @type {import('./$types').SubmitFunction} */
	const handleUpdate = async () => {
		isUpdating = true;
		return async ({ result }) => {
			isUpdating = false;
			if (result.type === 'success' || result.type === 'redirect') {
				close();
			}
			await applyAction(result);
		};
	};

	/** @type {import('./$types').SubmitFunction} */
	const handleUpload = async () => {
		isUploading = true;
		return async ({ result }) => {
			isUploading = false;
			/** @type {any} */
			const res = result;
			if (result.type === 'success' || result.type === 'redirect') {
				user.thumbnail = res.data.thumbnail;
			}
			await applyAction(result);
		};
	};
</script>

<div class="hero-container">
	<div class="hero-logo">
		<img
			src={user.thumbnail ? user.thumbnail : Avatar}
			alt={`${user.first_name} ${user.last_name}`}
		/>
	</div>
	<h3 class="hero-subtitle subtitle">
		Name (First and Last): {`${user.first_name} ${user.last_name}`}
	</h3>
	{#if user.profile.phone_number}
		<h3 class="hero-subtitle">
			Phone: {user.profile.phone_number}
		</h3>
	{/if}

	{#if user.profile.github_link}
		<h3 class="hero-subtitle">
			GitHub: {user.profile.github_link}
		</h3>
	{/if}

	{#if user.profile.birth_date}
		<h3 class="hero-subtitle">
			Date of birth: {user.profile.birth_date}
		</h3>
	{/if}
	<div class="hero-buttons-container">
		<button class="button-dark" on:click={open}>Edit profile</button>
	</div>
</div>

{#if showModal}
	<Modal on:close={close}>
		<form
			class="content image"
			action="?/uploadImage"
			method="post"
			enctype="multipart/form-data"
			use:enhance={handleUpload}
		>
			<ImageInput avatar={user.thumbnail} fieldName="thumbnail" title="Select user image" />

			{#if !user.thumbnail}
				<div class="btn-wrapper">
					{#if isUploading}
						<SmallLoader width={30} message={'Uploading...'} />
					{:else}
						<button class="button-dark" type="submit">Upload image</button>
					{/if}
				</div>
			{:else}
				<input type="hidden" hidden name="thumbnail_url" value={user.thumbnail} required />
				<div class="btn-wrapper">
					{#if isUploading}
						<SmallLoader width={30} message={'Removing...'} />
					{:else}
						<button class="button-dark" formaction="?/deleteImage" type="submit">
							Remove image
						</button>
					{/if}
				</div>
			{/if}
		</form>
		<form class="content" action="?/updateUser" method="POST" use:enhance={handleUpdate}>
			<h1 class="step-title" style="text-align: center;">Update User</h1>
			{#if form?.success}
				<h4
					class="step-subtitle warning"
					in:receive={{ key: Math.floor(Math.random() * 100) }}
					out:send={{ key: Math.floor(Math.random() * 100) }}
				>
					To avoid corrupt data and inconsistencies in your thumbnail, ensure you click on the
					"Update" button below.
				</h4>
			{/if}
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

			<input type="hidden" hidden name="thumbnail" value={user.thumbnail} />

			<div class="input-box">
				<span class="label">First name:</span>
				<input
					class="input"
					type="text"
					name="first_name"
					value={user.first_name}
					placeholder="Your first name..."
				/>
			</div>
			<div class="input-box">
				<span class="label">Last name:</span>
				<input
					class="input"
					type="text"
					name="last_name"
					value={user.last_name}
					placeholder="Your last name..."
				/>
			</div>
			<div class="input-box">
				<span class="label">Phone number:</span>
				<input
					class="input"
					type="tel"
					name="phone_number"
					value={user.profile.phone_number ? user.profile.phone_number : ''}
					placeholder="Your phone number e.g +2348135703593..."
				/>
			</div>
			<div class="input-box">
				<span class="label">Birth date:</span>
				<input
					class="input"
					type="date"
					name="birth_date"
					value={user.profile.birth_date ? user.profile.birth_date : ''}
					placeholder="Your date of birth..."
				/>
			</div>
			<div class="input-box">
				<span class="label">GitHub Link:</span>
				<input
					class="input"
					type="url"
					name="github_link"
					value={user.profile.github_link ? user.profile.github_link : ''}
					placeholder="Your github link e.g https://github.com/Sirneij/..."
				/>
			</div>
			{#if isUpdating}
				<SmallLoader width={30} message={'Updating...'} />
			{:else}
				<button type="submit" class="button-dark">Update</button>
			{/if}
		</form>
	</Modal>
{/if}

<style>
	.hero-container .hero-subtitle:not(:last-of-type) {
		margin: 0 0 0 0;
	}

	.content.image {
		display: flex;
		align-items: center;
		justify-content: center;
	}
	@media (max-width: 680px) {
		.content.image {
			margin: 0 0 0;
		}
	}
	.content.image .btn-wrapper {
		margin-top: 2.5rem;
		margin-left: 1rem;
	}
	.content.image .btn-wrapper button {
		padding: 15px 18px;
	}
</style>
