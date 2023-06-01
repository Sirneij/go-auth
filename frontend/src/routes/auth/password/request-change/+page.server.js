import { BASE_API_URI } from '$lib/utils/constants';
import { formatError, isEmpty, isValidEmail } from '$lib/utils/helpers';
import { fail, redirect } from '@sveltejs/kit';

/** @type {import('./$types').Actions} */
export const actions = {
	default: async ({ fetch, request }) => {
		const formData = await request.formData();
		const email = String(formData.get('email'));

		// Some validations
		/** @type {Record<string, string>} */
		const fieldsError = {};
		if (!isValidEmail(email)) {
			fieldsError.email = 'That email address is invalid.';
		}

		if (!isEmpty(fieldsError)) {
			return fail(400, { fieldsError: fieldsError });
		}
		/** @type {RequestInit} */
		const requestInitOptions = {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify({ email: email })
		};

		const res = await fetch(
			`${BASE_API_URI}/users/password/request-password-change/`,
			requestInitOptions
		);

		if (!res.ok) {
			const response = await res.json();
			const errors = formatError(response.error);
			return fail(400, { errors: errors });
		}

		const response = await res.json();

		// redirect the user
		throw redirect(302, `/auth/confirming?message=${response.message}`);
	}
};
