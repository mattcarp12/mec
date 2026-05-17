// --- lib/pb.ts ---
import PocketBase from 'pocketbase';

const pbUrl = process.env.NEXT_PUBLIC_PB_URL || 'http://127.0.0.1:8090';
export const pb = new PocketBase(pbUrl);