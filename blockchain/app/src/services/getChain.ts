import { DEFAULT_BASE_URL } from '@/lib/constants';
import type { GetChainResponse, Commonresponse } from '@/lib/types';

export async function getChain(
	baseUrl: string = DEFAULT_BASE_URL,
): Promise<GetChainResponse & Partial<Commonresponse>> {
	try {
		const res = await fetch(`${baseUrl}/chain`);
		if (!res.ok) {
			const errorData = await res.json().catch(() => ({}));
			throw new Error(errorData.error || 'Failed to get chain');
		}
		return await res.json();
	} catch (err: any) {
		return {
			chain: [],
			chain_length: 0,
			transaction_pool: [],
			blockchain_address: '',
			port: 0,
			mining_difficulty: 0,
			mining_reward: 0,
			mining: false,
			host: '',
			neighbors: [],
			wallets: {},
			message: 'failed',
			error: err.message || 'Unknown error occurred',
		};
	}
}
