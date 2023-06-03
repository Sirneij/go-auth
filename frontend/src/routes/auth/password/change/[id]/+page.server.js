import { BASE_API_URI } from '$lib/utils/constants';
import { formatError, isEmpty, isValidPasswordMedium } from '$lib/utils/helpers';
import { fail, redirect } from '@sveltejs/kit';

/** @type {import('./$types').Actions} */
export const actions = {
	default: async ({ fetch, request }) => {
		const formData = await request.formData();
		const password = String(formData.get('password'));
		const confirmPassword = String(formData.get('confirm_password'));
		let token = String(formData.get('token'));
		token = token.split(' ').join('');
		const userID = String(formData.get('user_id'));

		// Some validations
		/** @type {Record<string, string>} */
		const fieldsError = {};
		if (!isValidPasswordMedium(password)) {
			fieldsError.password =
				'Password is not valid. Password must contain six characters or more and has at least one lowercase and one uppercase alphabetical character or has at least one lowercase and one numeric character or has at least one uppercase and one numeric character.';
		}
		if (confirmPassword.trim() !== password.trim()) {
			fieldsError.confirmPassword = 'Password and confirm password do not match.';
		}

		if (!isEmpty(fieldsError)) {
			return fail(400, { fieldsError: fieldsError });
		}

		/** @type {RequestInit} */
		const requestInitOptions = {
			method: 'PUT',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify({ token: token, password: password })
		};

		const res = await fetch(
			`${BASE_API_URI}/users/password/change-user-password/${userID}/`,
			requestInitOptions
		);

		if (!res.ok) {
			const response = await res.json();
			const errors = formatError(response.error);
			return fail(400, { errors: errors });
		}

		const response = await res.json();

		// redirect the user
		throw redirect(302, `/auth/login?message=${response.message}`);
	}
};
