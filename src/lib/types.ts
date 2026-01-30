export interface Peer {
  id: string;
  name: string;
  publicKey: string;
  status: 'online' | 'offline';
  allowedIps: string[];
  latestHandshake: string;
  transfer: {
    received: number;
    sent: number;
  }
}
