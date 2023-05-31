import { BASE_API_URI } from '$lib/utils/constants';
import { fail, redirect } from '@sveltejs/kit';

/** @type {import('./$types').PageServerLoad} */
export async function load({ locals }) {
	// redirect user if logged in
	if (locals.user) {
		throw redirect(302, '/');
	}
}

/** @type {import('./$types').Actions} */
export const actions = {
	/**
	 *
	 * @param request - The request object
	 * @param fetch - Fetch object from sveltekit
	 * @param cookies - SvelteKit's cookie object
	 * @returns Error data or redirects user to the home page or the previous page
	 */
	login: async ({ request, fetch, cookies }) => {
		const data = await request.formData();
		const email = String(data.get('email'));
		const password = String(data.get('password'));
		const next = String(data.get('next'));

		/** @type {RequestInit} */
		const requestInitOptions = {
			method: 'POST',
			credentials: 'include',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify({
				email: email,
				password: password
			})
		};

		const res = await fetch(`${BASE_API_URI}/users/login/`, requestInitOptions);

		if (!res.ok) {
			const response = await res.json();
			const errors = [];
			errors.push({
				error: response.error.charAt(0).toUpperCase() + response.error.slice(1),
				id: 0
			});
			return fail(400, { errors: errors });
		}

		if (res.headers.has('Set-Cookie')) {
			const sessionID = Object.fromEntries(res.headers)
				['set-cookie'].split(';')[0]
				.split(/=(.*)/s)[1];

			cookies.set('go-auth-sessionid', sessionID, {
				httpOnly: true,
				sameSite: 'lax',
				path: '/',
				secure: true
			});
		}

		throw redirect(303, next || '/');
	}
};
