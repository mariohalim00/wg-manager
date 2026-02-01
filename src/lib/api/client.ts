// Base HTTP client for API communication
import type { APIResponse, APIError } from '../types/api';

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';

/**
 * Base fetch wrapper with error handling
 */
async function request<T>(endpoint: string, options: RequestInit = {}): Promise<APIResponse<T>> {
	const url = `${API_BASE_URL}${endpoint}`;

	try {
		const response = await fetch(url, {
			...options,
			headers: {
				'Content-Type': 'application/json',
				...options.headers
			}
		});

		const status = response.status;

		// Handle successful responses
		if (response.ok) {
			// 204 No Content - successful with no body
			if (status === 204) {
				return { status, data: undefined as T };
			}

			const data = await response.json();
			return { status, data };
		}

		// Handle error responses
		let error: APIError;
		try {
			const errorData = await response.json();
			error = {
				error: errorData.error || `HTTP ${status} error`,
				details: errorData.details
			};
		} catch {
			error = {
				error: `HTTP ${status} error`,
				details: response.statusText
			};
		}

		return { status, error };
	} catch (err) {
		// Network errors, timeouts, etc.
		const error: APIError = {
			error: 'Network error',
			details: err instanceof Error ? err.message : 'Failed to connect to backend'
		};
		return { status: 0, error };
	}
}

/**
 * GET request
 */
export async function get<T>(endpoint: string): Promise<APIResponse<T>> {
	return request<T>(endpoint, { method: 'GET' });
}

/**
 * POST request
 */
export async function post<T>(endpoint: string, body: unknown): Promise<APIResponse<T>> {
	return request<T>(endpoint, {
		method: 'POST',
		body: JSON.stringify(body)
	});
}

/**
 * DELETE request
 */
export async function del<T>(endpoint: string): Promise<APIResponse<T>> {
	return request<T>(endpoint, { method: 'DELETE' });
}

/**
 * PATCH request
 */
export async function patch<T>(endpoint: string, body: unknown): Promise<APIResponse<T>> {
	return request<T>(endpoint, {
		method: 'PATCH',
		body: JSON.stringify(body)
	});
}
