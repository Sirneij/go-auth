import { BASE_API_URI } from '$lib/utils/constants';
import { redirect } from '@sveltejs/kit';

/** @type {import('./$types').PageServerLoad} */
export async function load({ locals, cookies }) {
	// redirect user if not logged in or not a superuser
	if (!locals.user || !locals.user.is_superuser) {
		throw redirect(302, `/auth/login?next=/auth/logout`);
	}

	const fetchMetrics = async () => {
		const res = await fetch(`${BASE_API_URI}/metrics/`, {
			credentials: 'include',
			headers: {
				Cookie: `sessionid=${cookies.get('go-auth-sessionid')}`
			}
		});
		return res.ok && (await res.json());
	};

	return {
		metrics: fetchMetrics()
	};
}
