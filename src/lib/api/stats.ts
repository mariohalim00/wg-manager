// Stats API client
import { get } from './client';
import type { APIResponse } from '../types/api';
import type { InterfaceStats } from '../types/stats';

/**
 * Get interface statistics
 * GET /stats
 */
export async function getStats(): Promise<APIResponse<InterfaceStats>> {
	return get<InterfaceStats>('/stats');
}
