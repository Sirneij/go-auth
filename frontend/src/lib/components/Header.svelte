<script>
	import { applyAction, enhance } from '$app/forms';
	import { page } from '$app/stores';
	import Developer from '$lib/img/hero-image.png';
	import Avatar from '$lib/img/teamavatar.png';
</script>

<header class="header">
	<div class="header-container">
		<div class="header-left">
			<div class="header-crafted-by-container">
				<a href="https://github.com/Sirneij">
					<span>Developed by</span><img src={Developer} alt="John Owolabi Idogun" />
				</a>
			</div>
		</div>
		<div class="header-right">
			<div class="header-nav-item" class:active={$page.url.pathname === '/'}>
				<a href="/">home</a>
			</div>
			{#if !$page.data.user}
				<div class="header-nav-item" class:active={$page.url.pathname === '/auth/login'}>
					<a href="/auth/login">login</a>
				</div>
				<div class="header-nav-item" class:active={$page.url.pathname === '/auth/register'}>
					<a href="/auth/register">register</a>
				</div>
			{:else}
				<div class="header-nav-item">
					<a href="/auth/about/{$page.data.user.id}">
						<img
							src={$page.data.user.thumbnail ? $page.data.user.thumbnail : Avatar}
							alt={`${$page.data.user.first_name} ${$page.data.user.last_name}`}
						/>
					</a>
				</div>
				<form
					class="header-nav-item"
					action="/auth/logout"
					method="POST"
					use:enhance={async () => {
						return async ({ result }) => {
							await applyAction(result);
						};
					}}
				>
					<button type="submit">logout</button>
				</form>
			{/if}
		</div>
	</div>
</header>
