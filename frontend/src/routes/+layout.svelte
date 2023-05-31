<script>
	import Footer from '$lib/components/Footer.svelte';
	import Header from '$lib/components/Header.svelte';
	import Transition from '$lib/components/Transition.svelte';
	import '$lib/css/styles.min.css';

	/** @type {import('./$types').PageData} */
	export let data;
</script>

<Transition key={data.url} duration={600}>
	<Header />

	<slot />

	<Footer />
</Transition>
