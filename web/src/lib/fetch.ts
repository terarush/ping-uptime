import axios, { type AxiosInstance } from "axios";
import Cookies from 'js-cookie';

const Fetch = axios.create({
	baseURL: "/api",
	headers: {
		"Content-Type": "application/json",
	},
	withCredentials: true,
}) as AxiosInstance;


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

let isRefreshing = false;
let failedQueue: any[] = [];

const processQueue = (error: any, token: string | null = null) => {
	failedQueue.forEach((prom) => {
		if (error) {
			prom.reject(error);
		} else {
			prom.resolve(token);
		}
	});
	failedQueue = [];
};

Fetch.interceptors.response.use(
	(response) => {
		return response;
	},
	async (error) => {
		const originalRequest = error.config;

		// Catch 401 Unauthorized and auto-refresh the token
		if (error.response?.status === 401 && !originalRequest._retry) {
			const isAuthPath = originalRequest.url === "/auth/refresh" || originalRequest.url === "/auth/login";
			if (isAuthPath) {
				return Promise.reject(error);
			}

			if (isRefreshing) {
				return new Promise((resolve, reject) => {
					failedQueue.push({ resolve, reject });
				})
					.then((token) => {
						originalRequest.headers.Authorization = `Bearer ${token}`;
						return Fetch(originalRequest);
					})
					.catch((err) => {
						return Promise.reject(err);
					});
			}

			originalRequest._retry = true;
			isRefreshing = true;

			const refreshToken = Cookies.get("refreshToken");
			if (!refreshToken) {
				cleanupSessionAndRedirect();
				isRefreshing = false;
				return Promise.reject(error);
			}

			try {
				const response = await axios.post("/api/auth/refresh", {
					refreshToken: refreshToken,
				});

				const newAccessToken = response.data?.data?.accessToken;
				if (newAccessToken) {
					Cookies.set("accessToken", newAccessToken, { expires: 7 });
					Fetch.defaults.headers.common["Authorization"] = `Bearer ${newAccessToken}`;
					originalRequest.headers.Authorization = `Bearer ${newAccessToken}`;

					processQueue(null, newAccessToken);
					isRefreshing = false;
					return Fetch(originalRequest);
				}
			} catch (refreshError) {
				processQueue(refreshError, null);
				cleanupSessionAndRedirect();
				isRefreshing = false;
				return Promise.reject(refreshError);
			}
		}

		return Promise.reject(error);
	},
);

function cleanupSessionAndRedirect() {
	Cookies.remove("accessToken");
	Cookies.remove("refreshToken");
	if (typeof window !== "undefined") {
		window.location.href = "/";
	}
}

// Extend Fetch instance with uploadForm method for multipart/form-data
type ExtendedAxiosInstance = AxiosInstance & {
	uploadForm: <T = unknown>(url: string, data: FormData) => Promise<T>;
};

const ExtendedFetch = Fetch as ExtendedAxiosInstance;

ExtendedFetch.uploadForm = async <T = unknown>(
	url: string,
	data: FormData,
): Promise<T> => {
	const response = await Fetch.post(url, data);
	return response.data;
};

export default ExtendedFetch;
