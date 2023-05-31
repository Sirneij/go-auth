<script>
	import Footer from '$lib/components/Footer.svelte';
	import Header from '$lib/components/Header.svelte';
	import '$lib/css/styles.min.css';
</script>

<Header />

<main>
	<slot />
</main>

<Footer />
