import { DEFAULT_BASE_URL } from '@/lib/constants';
import type { GetWalletAmountResponse } from '@/lib/types';

export async function getWalletAmount(
	blockchain_address: string,
	baseUrl: string = DEFAULT_BASE_URL,
): Promise<GetWalletAmountResponse> {
	try {
		const res = await fetch(
			`${baseUrl}/address/${blockchain_address}/amount`,
		);

		if (!res.ok) {
			const errorData = await res.json().catch(() => ({}));
			throw new Error(errorData.error || `Failed to fetch wallet amount`);
		}

		return await res.json();
	} catch (err: any) {
		return {
			amount: 0,
			message: 'failed',
			error: err.message || 'Unknown error occurred',
		} as GetWalletAmountResponse;
	}
}
