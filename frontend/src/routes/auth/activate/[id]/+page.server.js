import { BASE_API_URI } from '$lib/utils/constants';
import { formatError } from '$lib/utils/helpers';
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
	 * @returns Error data or redirects user to the home page or the previous page
	 */
	default: async ({ request }) => {
		const data = await request.formData();
		const userID = String(data.get('user_id'));
		const token = String(data.get('token'));

		/** @type {RequestInit} */
		const requestInitOptions = {
			method: 'PUT',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify({
				token: token
			})
		};

		const res = await fetch(`${BASE_API_URI}/users/activate/${userID}/`, requestInitOptions);

		if (!res.ok) {
			const response = await res.json();
			const errors = formatError(response.error);

			return fail(400, { errors: errors });
		}

		const response = await res.json();

		throw redirect(303, `/auth/login?message=${response.message}`);
	}
};
