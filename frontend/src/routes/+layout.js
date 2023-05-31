/** @type {import('./$types').LayoutLoad} */
export async function load({ fetch, url, data }) {
	const { user } = data;
	return { fetch, url: url.pathname, user };
}
