// Validation utilities

/**
 * CIDR notation regex pattern
 * Matches IPv4 CIDR like 10.0.0.1/32 or 192.168.0.0/24
 */
const CIDR_REGEX =
	/^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\/(?:3[0-2]|[12]?[0-9])$/;

/**
 * Validate CIDR notation for allowed IPs
 * @param cidr - CIDR string to validate (e.g., "10.0.0.2/32")
 * @returns true if valid CIDR, false otherwise
 */
export function validateCIDR(cidr: string): boolean {
	return CIDR_REGEX.test(cidr.trim());
}

/**
 * Validate multiple CIDR strings (comma or newline separated)
 * @param input - String containing one or more CIDR notations
 * @returns Object with valid status and error message
 */
export function validateCIDRList(input: string): { valid: boolean; error?: string } {
	const cidrs = input
		.split(/[,\n]/)
		.map((s) => s.trim())
		.filter((s) => s.length > 0);

	if (cidrs.length === 0) {
		return { valid: false, error: 'At least one CIDR is required' };
	}

	const invalid = cidrs.filter((cidr) => !validateCIDR(cidr));

	if (invalid.length > 0) {
		return {
			valid: false,
			error: `Invalid CIDR notation: ${invalid.join(', ')}`
		};
	}

	return { valid: true };
}
