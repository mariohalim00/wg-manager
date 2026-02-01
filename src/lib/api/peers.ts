// Peer API client
import { get, post, del, patch } from './client';
import type { APIResponse } from '../types/api';
import type { Peer, PeerFormData, PeerCreateResponse, PeerUpdateRequest } from '../types/peer';

/**
 * List all peers
 * GET /peers
 */
export async function listPeers(): Promise<APIResponse<Peer[]>> {
	return get<Peer[]>('/peers');
}

/**
 * Add a new peer
 * POST /peers
 */
export async function addPeer(data: PeerFormData): Promise<APIResponse<PeerCreateResponse>> {
	return post<PeerCreateResponse>('/peers', data);
}

/**
 * Remove a peer by ID (public key)
 * DELETE /peers/{id}
 */
export async function removePeer(peerId: string): Promise<APIResponse<void>> {
	return del<void>(`/peers/${encodeURIComponent(peerId)}`);
}

/**
 * Update a peer's metadata or config
 * PATCH /peers/{id}
 */
export async function updatePeer(
	peerId: string,
	data: PeerUpdateRequest
): Promise<APIResponse<Peer>> {
	return patch<Peer>(`/peers/${encodeURIComponent(peerId)}`, data);
}

/**
 * Regenerate keys for a peer
 * POST /peers/regenerate-keys/{id}
 */
export async function regenerateKeys(peerId: string): Promise<APIResponse<PeerCreateResponse>> {
	return post<PeerCreateResponse>(`/peers/regenerate-keys/${encodeURIComponent(peerId)}`, {});
}
