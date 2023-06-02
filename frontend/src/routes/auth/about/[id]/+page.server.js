import { BASE_API_URI } from '$lib/utils/constants';
import { formatError } from '$lib/utils/helpers';
import { fail, redirect } from '@sveltejs/kit';

/** @type {import('./$types').PageServerLoad} */
export async function load({ locals, params }) {
	// redirect user if not logged in
	if (!locals.user) {
		throw redirect(302, `/auth/login?next=/auth/about/${params.id}`);
	}
}

/** @type {import('./$types').Actions} */
export const actions = {
	/**
	 *
	 * @param request - The request object
	 * @param fetch - Fetch object from sveltekit
	 * @param cookies - SvelteKit's cookie object
	 * @param locals - The local object, housing current user
	 * @returns Error data or redirects user to the home page or the previous page
	 */
	default: async ({ request, fetch, cookies, locals }) => {
		const formData = await request.formData();

		// Ensure that first_name is different from the current one
		if (formData.get('first_name')) {
			const firstName = formData.get('first_name');
			if (firstName === locals.user.first_name || firstName === '') {
				formData.delete('first_name');
			}
		}
		// Ensure that last_name is different from the current one
		if (formData.get('last_name')) {
			const lastName = formData.get('last_name');
			if (lastName === locals.user.last_name || lastName === '') {
				formData.delete('last_name');
			}
		}

		const apiURL = `${BASE_API_URI}/users/update-user/`;

		const res = await fetch(apiURL, {
			method: 'PATCH',
			credentials: 'include',
			headers: {
				'Content-Type': 'application/json',
				Cookie: `sessionid=${cookies.get('go-auth-sessionid')}`
			},
			body: formData
		});

		if (!res.ok) {
			const response = await res.json();
			const errors = formatError(response.error);
			return fail(400, { errors: errors });
		}

		const response = await res.json();

		locals.user = response;

		throw redirect(303, `/auth/about/${response.id}`);
	}
};
