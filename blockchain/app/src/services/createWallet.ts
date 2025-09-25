import { DEFAULT_BASE_URL } from '@/lib/constants';
import type { Wallet, Commonresponse } from '@/lib/types';

export async function createWallet(
	baseUrl: string = DEFAULT_BASE_URL,
): Promise<Wallet & Partial<Commonresponse>> {
	try {
		const res = await fetch(`${baseUrl}/wallet`);
		if (!res.ok) {
			const errorData = await res.json().catch(() => ({}));
			throw new Error(errorData.error || 'Failed to create wallet');
		}
		return await res.json();
	} catch (err: any) {
		return {
			private_key: '',
			public_key: '',
			blockchain_address: '',
			message: 'failed',
			error: err.message || 'Unknown error occurred',
		};
	}
}
