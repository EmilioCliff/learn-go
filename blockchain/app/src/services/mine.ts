import { DEFAULT_BASE_URL } from '@/lib/constants';
import type { Commonresponse } from '@/lib/types';

export async function mine(
	baseUrl: string = DEFAULT_BASE_URL,
): Promise<Commonresponse> {
	try {
		const res = await fetch(`${baseUrl}/mine`);
		if (!res.ok) {
			const errorData = await res.json().catch(() => ({}));
			throw new Error(errorData.error || 'Failed to mine');
		}
		return await res.json();
	} catch (err: any) {
		return {
			message: 'failed',
			error: err.message || 'Unknown error occurred',
		} as Commonresponse;
	}
}
