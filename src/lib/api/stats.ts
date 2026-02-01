// Stats API client
import { get } from './client';
import type { APIResponse } from '../types/api';
import type { InterfaceStats, StatsHistoryItem } from '../types/stats';

/**
 * Get interface statistics
 * GET /stats
 */
export async function getStats(): Promise<APIResponse<InterfaceStats>> {
	return get<InterfaceStats>('/stats');
}

/**
 * Get historical statistics
 * GET /stats/history
 */
export async function getStatsHistory(): Promise<APIResponse<StatsHistoryItem[]>> {
	return get<StatsHistoryItem[]>('/stats/history');
}
