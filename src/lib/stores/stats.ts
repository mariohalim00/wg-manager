import { writable } from 'svelte/store';

// Define a type for your stats for better type safety
export interface Stats {
  totalPeers: number;
  onlinePeers: number;
  totalDataUsage: {
    sent: number;
    received: number;
  };
}

export const stats = writable<Stats>({
  totalPeers: 0,
  onlinePeers: 0,
  totalDataUsage: {
    sent: 0,
    received: 0,
  }
});
