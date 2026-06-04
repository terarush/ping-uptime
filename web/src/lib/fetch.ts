import axios, { type AxiosInstance } from "axios";

const Fetch = axios.create({
	baseURL: `${import.meta.env.VITE_API_URL}/api/v1`,
	headers: {
		"Content-Type": "application/json",
	},
	withCredentials: true,
}) as AxiosInstance;

import Cookies from 'js-cookie';

Fetch.interceptors.request.use(
	(config) => {
		const token = Cookies.get("accessToken");
		if (token && config.headers) {
			config.headers.Authorization = `Bearer ${token}`;
		}
		return config;
	},
	(error) => {
		return Promise.reject(error);
	},
);

// Extend Fetch instance with uploadForm method for multipart/form-data
type ExtendedAxiosInstance = AxiosInstance & {
	uploadForm: <T = unknown>(url: string, data: FormData) => Promise<T>;
};

const ExtendedFetch = Fetch as ExtendedAxiosInstance;

ExtendedFetch.uploadForm = async <T = unknown>(
	url: string,
	data: FormData,
): Promise<T> => {
	const response = await Fetch.post(url, data, {
		headers: {
			"Content-Type": "multipart/form-data",
		},
	});
	return response.data;
};

export default ExtendedFetch;
