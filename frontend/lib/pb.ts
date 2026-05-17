// --- lib/pb.ts ---
import PocketBase from 'pocketbase';

// Connect to your local PocketBase container
export const pb = new PocketBase('http://127.0.0.1:8090');