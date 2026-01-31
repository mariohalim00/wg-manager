// Formatting utilities

/**
 * Format bytes to human-readable string (B, KB, MB, GB, TB)
 * @param bytes - Number of bytes to format
 * @returns Formatted string with appropriate unit
 */
export function formatBytes(bytes: number): string {
	if (bytes === 0) return '0 B';

	const units = ['B', 'KB', 'MB', 'GB', 'TB'];
	const k = 1024;
	const i = Math.floor(Math.log(bytes) / Math.log(k));
	const value = bytes / Math.pow(k, i);

	// Format to 2 decimal places, but remove trailing zeros
	const formatted = value.toFixed(2).replace(/\.?0+$/, '');

	return `${formatted} ${units[i]}`;
}

/**
 * Format last handshake time to relative string or "Never"
 * @param lastHandshake - ISO timestamp string or "0"
 * @returns Relative time string (e.g., "2 minutes ago") or "Never"
 */
export function formatLastHandshake(lastHandshake: string): string {
	if (lastHandshake === '0' || !lastHandshake) {
		return 'Never';
	}

	const handshakeTime = new Date(lastHandshake).getTime();
	const now = Date.now();
	const diff = now - handshakeTime;

	// Less than 1 minute
	if (diff < 60000) {
		const seconds = Math.floor(diff / 1000);
		return seconds === 1 ? '1 second ago' : `${seconds} seconds ago`;
	}

	// Less than 1 hour
	if (diff < 3600000) {
		const minutes = Math.floor(diff / 60000);
		return minutes === 1 ? '1 minute ago' : `${minutes} minutes ago`;
	}

	// Less than 1 day
	if (diff < 86400000) {
		const hours = Math.floor(diff / 3600000);
		return hours === 1 ? '1 hour ago' : `${hours} hours ago`;
	}

	// Days
	const days = Math.floor(diff / 86400000);
	return days === 1 ? '1 day ago' : `${days} days ago`;
}

/**
 * Format date to local string
 * @param dateString - ISO timestamp string
 * @returns Formatted local date string
 */
export function formatDate(dateString: string): string {
	if (!dateString || dateString === '0') {
		return 'N/A';
	}

	const date = new Date(dateString);
	return date.toLocaleString();
}
