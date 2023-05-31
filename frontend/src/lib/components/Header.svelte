<script>
	import { page } from '$app/stores';
	import Developer from '$lib/img/hero-image.png';
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
				<a href="/">Home</a>
			</div>
			{#if !$page.data.user}
				<div class="header-nav-item" class:active={$page.url.pathname === '/auth/login'}>
					<a href="/auth/login">Login</a>
				</div>
			{:else}
				<div class="header-nav-item" class:active={$page.url.pathname === '/auth/logout'}>
					<a href="/auth/logout">Logout</a>
				</div>
			{/if}
		</div>
	</div>
</header>
