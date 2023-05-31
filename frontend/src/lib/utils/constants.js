export const BASE_API_URI = import.meta.env.DEV
	? import.meta.env.VITE_BASE_API_URI_DEV
	: import.meta.env.VITE_BASE_API_URI_PROD;
export const BASE_SERIES_API_URI = import.meta.env.DEV
	? import.meta.env.VITE_BASE_SERIES_API_URI_DEV
	: import.meta.env.VITE_BASE_SERIES_API_URI_PROD;
