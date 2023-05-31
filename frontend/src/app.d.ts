// See https://kit.svelte.dev/docs/types#app
// for information about these interfaces

interface UserProfile {
	id: string;
	user_id: string;
	phone_number: string | null;
	birth_date: string | null;
	github_link: string | null;
}

interface User {
	email: string;
	first_name: string;
	last_name: string;
	id: string;
	is_staff: boolean;
	thumbnail: string;
	is_superuser: boolean;
	profile: UserProfile;
}
declare global {
	namespace App {
		// interface Error {}
		interface Locals {
			user: User;
		}
		// interface PageData {}
		// interface Platform {}
	}
}

export {};
