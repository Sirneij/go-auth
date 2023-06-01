// @ts-nocheck
import { quintOut } from 'svelte/easing';
import { crossfade } from 'svelte/transition';

export const [send, receive] = crossfade({
	duration: (d) => Math.sqrt(d * 200),

	// eslint-disable-next-line no-unused-vars
	fallback(node, params) {
		const style = getComputedStyle(node);
		const transform = style.transform === 'none' ? '' : style.transform;

		return {
			duration: 600,
			easing: quintOut,
			css: (t) => `
                transform: ${transform} scale(${t});
                opacity: ${t}
            `
		};
	}
});

/**
 * Validates an email field
 * @file lib/utils/helpers/input.validation.ts
 * @param {string} email - The email to validate
 */
export const isValidEmail = (email) => {
	const EMAIL_REGEX =
		/[a-z0-9!#$%&'*+/=?^_`{|}~-]+(?:\.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*@(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?/;
	return EMAIL_REGEX.test(email.trim());
};
/**
 * Validates a strong password field
 * @file lib/utils/helpers/input.validation.ts
 * @param {string} password - The password to validate
 */
export const isValidPasswordStrong = (password) => {
	const strongRegex = new RegExp('^(?=.*[a-z])(?=.*[A-Z])(?=.*[0-9])(?=.*[!@#$%^&*])(?=.{8,})');

	return strongRegex.test(password.trim());
};
/**
 * Validates a medium password field
 * @file lib/utils/helpers/input.validation.ts
 * @param {string} password - The password to validate
 */
export const isValidPasswordMedium = (password) => {
	const mediumRegex = new RegExp(
		'^(((?=.*[a-z])(?=.*[A-Z]))|((?=.*[a-z])(?=.*[0-9]))|((?=.*[A-Z])(?=.*[0-9])))(?=.{6,})'
	);

	return mediumRegex.test(password.trim());
};

/**
 * Test whether or not an object is empty.
 * @param {Record<string, string>} obj - The object to test
 * @returns `true` or `false`
 */

export function isEmpty(obj) {
	for (const _i in obj) {
		return false;
	}
	return true;
}
/**
 * Test whether or not an object is empty.
 * @param {any} obj - The object to test
 * @returns `true` or `false`
 */

export function formatError(obj) {
	const errors = [];
	if (typeof obj === 'object' && obj !== null) {
		if (Array.isArray(obj)) {
			obj.forEach((/** @type {Object} */ error) => {
				Object.keys(error).map((k) => {
					errors.push({
						error: error[k],
						id: Math.random() * 1000
					});
				});
			});
		} else {
			Object.keys(obj).map((k) => {
				errors.push({
					error: obj[k],
					id: Math.random() * 1000
				});
			});
		}
	} else {
		errors.push({
			error: obj.charAt(0).toUpperCase() + obj.slice(1),
			id: 0
		});
	}

	return errors;
}
