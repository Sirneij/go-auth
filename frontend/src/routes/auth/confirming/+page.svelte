<script>
	import { page } from '$app/stores';
	let message = '';
	if ($page.url.search) {
		message = $page.url.search.split('=')[1].replaceAll('%20', ' ');
	}
</script>

<div class="container">
	<div class="content">
		<h1 class="step-title title">Email sent!</h1>
		<h4 class="step-subtitle normal">
			{message}
		</h4>
	</div>
</div>
