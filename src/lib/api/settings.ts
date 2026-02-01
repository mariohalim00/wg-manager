// Settings API client
import { get, post } from './client';
import type { APIResponse } from '../types/api';
import type { GlobalSettings } from '../types/settings';

/**
 * Get global settings
 * GET /settings
 */
export async function getSettings(): Promise<APIResponse<GlobalSettings>> {
	return get<GlobalSettings>('/settings');
}

/**
 * Update global settings
 * POST /settings
 */
export async function updateSettings(settings: GlobalSettings): Promise<APIResponse<void>> {
	return post<void>('/settings', settings);
}
