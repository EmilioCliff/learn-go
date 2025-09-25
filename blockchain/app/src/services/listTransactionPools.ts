import { DEFAULT_BASE_URL } from '@/lib/constants';
import type { ListTransactionPoolsResponse, Commonresponse } from '@/lib/types';

export async function listTransactionPools(
	baseUrl: string = DEFAULT_BASE_URL,
): Promise<ListTransactionPoolsResponse & Partial<Commonresponse>> {
	try {
		const res = await fetch(`${baseUrl}/transactions`);
		if (!res.ok) {
			const errorData = await res.json().catch(() => ({}));
			throw new Error(
				errorData.error || 'Failed to list transaction pools',
			);
		}
		return await res.json();
	} catch (err: any) {
		return {
			message: 'failed',
			error: err.message || 'Unknown error occurred',
		} as ListTransactionPoolsResponse & Partial<Commonresponse>;
	}
}
