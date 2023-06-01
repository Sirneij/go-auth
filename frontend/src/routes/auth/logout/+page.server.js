import { BASE_API_URI } from '$lib/utils/constants';
import { fail, redirect } from '@sveltejs/kit';

/** @type {import('./$types').PageServerLoad} */
export async function load({ locals }) {
	// redirect user if not logged in
	if (!locals.user) {
		throw redirect(302, `/auth/login?next=/auth/logout`);
	}
}

/** @type {import('./$types').Actions} */
export const actions = {
	default: async ({ fetch, cookies }) => {
		/** @type {RequestInit} */
		const requestInitOptions = {
			method: 'POST',
			credentials: 'include',
			headers: {
				'Content-Type': 'application/json',
				Cookie: `sessionid=${cookies.get('go-auth-sessionid')}`
			}
		};

		const res = await fetch(`${BASE_API_URI}/users/logout/`, requestInitOptions);

		if (!res.ok) {
			const response = await res.json();
			const errors = [];
			errors.push({ error: response.error, id: 0 });
			return fail(400, { errors: errors });
		}

		// eat the cookie
		cookies.delete('go-auth-sessionid', { path: '/' });

		// redirect the user
		throw redirect(302, '/auth/login');
	}
};
