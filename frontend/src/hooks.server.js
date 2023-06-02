import { BASE_API_URI } from '$lib/utils/constants';

/** @type {import('@sveltejs/kit').Handle} */
export async function handle({ event, resolve }) {
	if (event.locals.user) {
		// if there is already a user  in session load page as normal
		return await resolve(event);
	}
	// get cookies from browser
	const session = event.cookies.get('go-auth-sessionid');

	if (!session) {
		// if there is no session load page as normal
		return await resolve(event);
	}

	// find the user based on the session
	const res = await event.fetch(`${BASE_API_URI}/users/current-user/`, {
		credentials: 'include',
		headers: {
			Cookie: `sessionid=${session}`
		}
	});

	if (!res.ok) {
		// if there is no session load page as normal
		return await resolve(event);
	}

	// if `user` exists set `events.local`
	const response = await res.json();

	event.locals.user = response;
	if (event.locals.user.profile.birth_date) {
		event.locals.user.profile.birth_date = response['profile']['birth_date'].split('T')[0];
	}

	// load page as normal
	return await resolve(event);
}
