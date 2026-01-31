// Application constants

/**
 * API base URL from environment variable
 */
export const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';

/**
 * Handshake timeout threshold in milliseconds
 * Peer is considered online if lastHandshake is within this time
 * FR-002: 120 seconds = 120000 milliseconds
 */
export const HANDSHAKE_TIMEOUT = 120000;

/**
 * Default notification durations by type (milliseconds)
 */
export const NOTIFICATION_DURATIONS = {
	success: 3000,
	error: 5000,
	warning: 4000,
	info: 3000
} as const;

/**
 * Responsive breakpoints (matches Tailwind defaults)
 */
export const BREAKPOINTS = {
	sm: 640,
	md: 768,
	lg: 1024,
	xl: 1280,
	'2xl': 1536
} as const;

/**
 * Data refresh intervals (milliseconds)
 */
export const REFRESH_INTERVALS = {
	peers: 5000, // 5 seconds
	stats: 10000 // 10 seconds
} as const;
