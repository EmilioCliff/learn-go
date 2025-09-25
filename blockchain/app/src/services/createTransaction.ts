import type { Commonresponse } from '@/lib/types';

interface CreateTransactionParams {
	transaction: {
		sender_private_key: string;
		sender_public_key: string;
		sender_blockchain_address: string;
		recipient_blockchain_address: string;
		value: number;
	};
	baseUrl?: string;
}

export async function createTransaction(
	data: CreateTransactionParams,
): Promise<Commonresponse> {
	try {
		const res = await fetch(`${data.baseUrl}/transactions`, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json',
				Accept: 'application/json',
			},
			body: JSON.stringify(data.transaction),
		});
		if (!res.ok) {
			const errorData = await res.json().catch(() => ({}));
			throw new Error(errorData.error || 'Failed to create transaction');
		}
		return await res.json();
	} catch (error: any) {
		if (error.response) {
			throw new Error(error.response.error);
		}

		throw new Error(error.message);
	}
}
