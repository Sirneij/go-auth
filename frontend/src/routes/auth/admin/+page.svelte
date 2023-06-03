<script>
	import { page } from '$app/stores';
	/** @type {import('./$types').PageData} */
	export let data;

	$: ({ metrics } = data);

	const calculateAvgProTime = (/** @type {any} */ metric) => {
		const div = metric.total_processing_time_μs / metric.total_requests_received;
		const inSecs = div * 0.000001;
		return `${inSecs.toFixed(2)}s`;
	};
</script>

<div class="app">
	<div class="app-body">
		<nav class="navigation">
			<a href="/auth/admin" class:active={$page.url.pathname === '/auth/admin'}>
				<svg xmlns="http://www.w3.org/2000/svg" height="1em" viewBox="0 0 512 512">
					<path
						d="M0 256a256 256 0 1 1 512 0A256 256 0 1 1 0 256zm320 96c0-26.9-16.5-49.9-40-59.3V88c0-13.3-10.7-24-24-24s-24 10.7-24 24V292.7c-23.5 9.5-40 32.5-40 59.3c0 35.3 28.7 64 64 64s64-28.7 64-64zM144 176a32 32 0 1 0 0-64 32 32 0 1 0 0 64zm-16 80a32 32 0 1 0 -64 0 32 32 0 1 0 64 0zm288 32a32 32 0 1 0 0-64 32 32 0 1 0 0 64zM400 144a32 32 0 1 0 -64 0 32 32 0 1 0 64 0z"
					/>
				</svg>
				<span>Metrics</span>
			</a>
			<a href="/auth/admin#" class:active={$page.url.pathname === '/auth/admin#'}>
				<svg xmlns="http://www.w3.org/2000/svg" height="1em" viewBox="0 0 640 512">
					<path
						d="M144 0a80 80 0 1 1 0 160A80 80 0 1 1 144 0zM512 0a80 80 0 1 1 0 160A80 80 0 1 1 512 0zM0 298.7C0 239.8 47.8 192 106.7 192h42.7c15.9 0 31 3.5 44.6 9.7c-1.3 7.2-1.9 14.7-1.9 22.3c0 38.2 16.8 72.5 43.3 96c-.2 0-.4 0-.7 0H21.3C9.6 320 0 310.4 0 298.7zM405.3 320c-.2 0-.4 0-.7 0c26.6-23.5 43.3-57.8 43.3-96c0-7.6-.7-15-1.9-22.3c13.6-6.3 28.7-9.7 44.6-9.7h42.7C592.2 192 640 239.8 640 298.7c0 11.8-9.6 21.3-21.3 21.3H405.3zM224 224a96 96 0 1 1 192 0 96 96 0 1 1 -192 0zM128 485.3C128 411.7 187.7 352 261.3 352H378.7C452.3 352 512 411.7 512 485.3c0 14.7-11.9 26.7-26.7 26.7H154.7c-14.7 0-26.7-11.9-26.7-26.7z"
					/>
				</svg>
				<span>Users</span>
			</a>
		</nav>

		<div class="app-body-main-content">
			<section class="service-section">
				<h2>Metrics ({metrics.version})</h2>
				<div class="tiles">
					<article class="tile">
						<div class="tile-header">
							<svg xmlns="http://www.w3.org/2000/svg" height="1em" viewBox="0 0 448 512">
								<path
									d="M349.4 44.6c5.9-13.7 1.5-29.7-10.6-38.5s-28.6-8-39.9 1.8l-256 224c-10 8.8-13.6 22.9-8.9 35.3S50.7 288 64 288H175.5L98.6 467.4c-5.9 13.7-1.5 29.7 10.6 38.5s28.6 8 39.9-1.8l256-224c10-8.8 13.6-22.9 8.9-35.3s-16.6-20.7-30-20.7H272.5L349.4 44.6z"
								/>
							</svg>
							<h3>
								<span>Avg Pro. Time</span>
								<span>total pro. time / total reqs</span>
							</h3>
						</div>
						<p>{calculateAvgProTime(metrics)}</p>
						<div>{`${metrics.total_processing_time_μs} / ${metrics.total_requests_received}`}</div>
					</article>
					<article class="tile">
						<div class="tile-header">
							<svg xmlns="http://www.w3.org/2000/svg" height="1em" viewBox="0 0 640 512">
								<path
									d="M256 0c-35 0-64 59.5-64 93.7v84.6L8.1 283.4c-5 2.8-8.1 8.2-8.1 13.9v65.5c0 10.6 10.2 18.3 20.4 15.4l171.6-49 0 70.9-57.6 43.2c-4 3-6.4 7.8-6.4 12.8v42c0 7.8 6.3 14 14 14c1.3 0 2.6-.2 3.9-.5L256 480l110.1 31.5c1.3 .4 2.6 .5 3.9 .5c6 0 11.1-3.7 13.1-9C344.5 470.7 320 422.2 320 368c0-60.6 30.6-114 77.1-145.6L320 178.3V93.7C320 59.5 292 0 256 0zM640 368a144 144 0 1 0 -288 0 144 144 0 1 0 288 0zm-76.7-43.3c6.2 6.2 6.2 16.4 0 22.6l-72 72c-6.2 6.2-16.4 6.2-22.6 0l-40-40c-6.2-6.2-6.2-16.4 0-22.6s16.4-6.2 22.6 0L480 385.4l60.7-60.7c6.2-6.2 16.4-6.2 22.6 0z"
								/>
							</svg>
							<h3>
								<span>Active in-flight reqs</span>
								<span>total reqs - total res</span>
							</h3>
						</div>
						<p>{metrics.total_requests_received - metrics.total_responses_sent}</p>
						<div>{`${metrics.total_requests_received} - ${metrics.total_responses_sent}`}</div>
					</article>
					<article class="tile">
						<div class="tile-header">
							<svg xmlns="http://www.w3.org/2000/svg" height="1em" viewBox="0 0 512 512">
								<path
									d="M448 160H320V128H448v32zM48 64C21.5 64 0 85.5 0 112v64c0 26.5 21.5 48 48 48H464c26.5 0 48-21.5 48-48V112c0-26.5-21.5-48-48-48H48zM448 352v32H192V352H448zM48 288c-26.5 0-48 21.5-48 48v64c0 26.5 21.5 48 48 48H464c26.5 0 48-21.5 48-48V336c0-26.5-21.5-48-48-48H48z"
								/>
							</svg>
							<h3>
								<span>Goroutines</span>
								<span>No. of active goroutines</span>
							</h3>
						</div>
						<p>{metrics.goroutines}</p>
						<div>No. of active goroutines</div>
					</article>
				</div>
			</section>
		</div>
	</div>
</div>

<style>
	a {
		color: inherit;
	}

	.app {
		min-height: 80vh;
		width: 90%;
		max-width: 95%;
		padding: 2vw 4vw 6vw;
		display: flex;
		flex-direction: column;
	}

	.app-body {
		height: 100%;
		display: flex;
		flex-direction: column;
		padding-top: 2.5rem;
	}
	@media (max-width: 1200px) {
		.app-body > * {
			margin-bottom: 3.5rem;
		}
	}

	.navigation {
		display: flex;
		flex-direction: row;
		align-items: baseline;
		gap: 1rem;
		margin-bottom: 2rem;
	}
	.navigation a {
		display: flex;
		align-items: center;
		text-decoration: none;
		transition: 0.25s ease;
		border: none;
	}

	.navigation a * {
		transition: 0.25s ease;
	}
	.navigation a svg {
		margin-right: 0.75rem;
		fill: #595b5c;
		font-size: 1.25em;
		flex-shrink: 0;
	}

	.navigation a + a {
		margin-top: 1.25rem;
	}
	.navigation a:hover,
	.navigation a:focus {
		border: none;
		transform: translateX(4px);
		color: white;
	}

	.navigation a.active {
		border: none;
		transform: translateX(4px);
		color: white;
	}

	.navigation a:hover svg,
	.navigation a:focus svg {
		fill: white;
	}
	.navigation a.active svg {
		fill: white;
	}

	.tiles {
		display: flex;
		overflow-y: hidden;
		overflow-x: auto;
	}
	@media (max-width: 700px) {
		.tiles {
			display: grid;
			grid-template-columns: repeat(1, 1fr);
			-moz-column-gap: 1rem;
			column-gap: 1rem;
			row-gap: 1rem;
			margin-top: 1.25rem;
		}
	}

	.tile {
		padding: 1rem;
		border-radius: 8px;
		background-color: #2f2f2f;
		min-height: 200px;
		max-width: 330px;
		display: flex;
		flex-direction: column;
		justify-content: space-between;
		position: relative;
		transition: 0.25s ease;
		margin-right: 2rem;
	}
	.tile:hover {
		transform: translateY(-5px);
	}
	.tile:focus-within {
		box-shadow: 0 0 0 2px rgba(0, 0, 0, 0.25), 0 0 0 4px rgba(0, 0, 0, 0.75);
	}

	.tile div {
		font-family: 'Roboto Mono', monospace;
		text-decoration: none;
		display: flex;
		align-items: center;
		justify-content: space-between;
		font-size: 0.8rem;
		font-weight: 200;
	}

	.tile-header {
		display: flex;
		align-items: center;
	}
	.tile-header svg {
		transition: fill 0.3s ease-in-out;
		fill: #d43aff;
		font-size: 2.5rem;
	}
	.tile:hover svg {
		fill: #c907ff;
	}
	.tile-header h3 {
		display: flex;
		flex-direction: column;
		line-height: 1.375;
		margin-left: 0.5rem;
	}
	.tile-header h3 span:first-child {
		font-weight: 600;
	}
	.tile-header h3 span:last-child {
		font-family: 'Roboto Mono', monospace;
		font-size: 0.8rem;
		font-weight: 200;
	}
	.tile p {
		transition: all 0.3s ease-in-out;
		font-size: 2.5rem;
		text-align: center;
		color: white;
	}

	.tile:hover p {
		color: whitesmoke;
	}

	.service-section > h2 {
		font-size: 1.5rem;
		margin-bottom: 1.25rem;
	}

	.payment-details div span {
		font-size: 1.5rem;
	}
</style>
