// API error and response types

export interface APIError {
	error: string; // User-facing error message
	details?: string; // Optional detailed error
}

export interface APIResponse<T> {
	data?: T;
	error?: APIError;
	status: number;
}
