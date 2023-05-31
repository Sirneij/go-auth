/** @type {import('./$types').LayoutLoad} */
export async function load({ fetch, url }) {
	return { fetch, url: url.pathname };
}
