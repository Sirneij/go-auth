<script>
	// @ts-nocheck
	export let avatar;

	export let fieldName;
	export let title;

	let newAvatar;
	const onFileSelected = (e) => {
		const target = e.target;
		if (target && target.files) {
			let reader = new FileReader();
			reader.readAsDataURL(target.files[0]);
			reader.onload = (e) => {
				newAvatar = e.target?.result;
			};
		}
	};
</script>

<div id="app">
	{#if avatar}
		<img class="avatar" src={avatar} alt="d" />
	{:else}
		<img
			class="avatar"
			src={newAvatar
				? newAvatar
				: 'https://cdn4.iconfinder.com/data/icons/small-n-flat/24/user-alt-512.png'}
			alt=""
		/>
		<input type="file" id="file" name={fieldName} required on:change={(e) => onFileSelected(e)} />
		<label for="file" class="btn-3">
			{#if newAvatar}
				<span>Image selected! Click upload.</span>
			{:else}
				<span>{title}</span>
			{/if}
		</label>
	{/if}
</div>

<style>
	#app {
		margin-top: 1rem;
		display: flex;
		align-items: center;
		justify-content: center;
		flex-flow: column;
		color: rgb(148 163 184);
	}

	.avatar {
		display: flex;
		height: 6.5rem;
		width: 8rem;
	}
	[type='file'] {
		height: 0;
		overflow: hidden;
		width: 0;
	}
	[type='file'] + label {
		background: #9b9b9b;
		border: none;
		border-radius: 5px;
		color: #fff;
		cursor: pointer;
		display: inline-block;
		font-weight: 500;
		margin-bottom: 1rem;
		outline: none;
		padding: 1rem 50px;
		position: relative;
		transition: all 0.3s;
		vertical-align: middle;
	}
	[type='file'] + label:hover {
		background-color: #9b9b9b;
	}
	[type='file'] + label.btn-3 {
		background-color: #d43aff;
		border-radius: 0;
		overflow: hidden;
	}
	[type='file'] + label.btn-3 span {
		display: inline-block;
		height: 100%;
		transition: all 0.3s;
		width: 100%;
	}
	[type='file'] + label.btn-3::before {
		color: #fff;
		content: '\01F4F7';
		font-size: 200%;
		height: 100%;
		left: 45%;
		position: absolute;
		top: -180%;
		transition: all 0.3s;
		width: 100%;
	}
	[type='file'] + label.btn-3:hover {
		background-color: rgba(14, 166, 236, 0.5);
	}
	[type='file'] + label.btn-3:hover span {
		transform: translateY(300%);
	}
	[type='file'] + label.btn-3:hover::before {
		top: 0;
	}
</style>
