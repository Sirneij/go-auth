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
	updateUser: async ({ request, fetch, cookies, locals }) => {
		const formData = await request.formData();
		const firstName = String(formData.get('first_name'));
		const lastName = String(formData.get('last_name'));
		const thumbnail = String(formData.get('thumbnail'));
		const phoneNumber = String(formData.get('phone_number'));
		const birthDate = String(formData.get('birth_date'));
		const githubLink = String(formData.get('github_link'));

		const apiURL = `${BASE_API_URI}/users/update-user/`;

		const res = await fetch(apiURL, {
			method: 'PATCH',
			credentials: 'include',
			headers: {
				'Content-Type': 'application/json',
				Cookie: `sessionid=${cookies.get('go-auth-sessionid')}`
			},
			body: JSON.stringify({
				first_name: firstName,
				last_name: lastName,
				thumbnail: thumbnail,
				phone_number: phoneNumber,
				birth_date: birthDate,
				github_link: githubLink
			})
		});

		if (!res.ok) {
			const response = await res.json();
			const errors = formatError(response.error);
			return fail(400, { errors: errors });
		}

		const response = await res.json();

		locals.user = response;

		if (locals.user.profile.birth_date) {
			locals.user.profile.birth_date = response['profile']['birth_date'].split('T')[0];
		}

		throw redirect(303, `/auth/about/${response.id}`);
	},
	/**
	 *
	 * @param request - The request object
	 * @param fetch - Fetch object from sveltekit
	 * @param cookies - SvelteKit's cookie object
	 * @param locals - The local object, housing current user
	 * @returns Error data or redirects user to the home page or the previous page
	 */
	uploadImage: async ({ request, fetch, cookies }) => {
		const formData = await request.formData();

		/** @type {RequestInit} */
		const requestInitOptions = {
			method: 'POST',
			headers: {
				Cookie: `sessionid=${cookies.get('go-auth-sessionid')}`
			},
			body: formData
		};

		const res = await fetch(`${BASE_API_URI}/file/upload/`, requestInitOptions);

		if (!res.ok) {
			const response = await res.json();
			const errors = formatError(response.error);
			return fail(400, { errors: errors });
		}

		const response = await res.json();

		return {
			success: true,
			thumbnail: response['s3_url']
		};
	},
	/**
	 *
	 * @param request - The request object
	 * @param fetch - Fetch object from sveltekit
	 * @param cookies - SvelteKit's cookie object
	 * @param locals - The local object, housing current user
	 * @returns Error data or redirects user to the home page or the previous page
	 */
	deleteImage: async ({ request, fetch, cookies }) => {
		const formData = await request.formData();

		/** @type {RequestInit} */
		const requestInitOptions = {
			method: 'DELETE',
			headers: {
				Cookie: `sessionid=${cookies.get('go-auth-sessionid')}`
			},
			body: formData
		};

		const res = await fetch(`${BASE_API_URI}/file/delete/`, requestInitOptions);

		if (!res.ok) {
			const response = await res.json();
			const errors = formatError(response.error);
			return fail(400, { errors: errors });
		}

		return {
			success: true,
			thumbnail: ''
		};
	}
};
