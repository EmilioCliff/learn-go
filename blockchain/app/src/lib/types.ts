export interface Commonresponse {
	message: string;
	error: string;
}

export interface Wallet {
	private_key: string;
	public_key: string;
	blockchain_address: string;
	amount?: number;
}

export interface Transaction {
	sender_blockchain_address: string;
	recipient_blockchain_address: string;
	value: number;
	timestamp?: number;
}

export interface Block {
	nonce: number;
	previous_hash: string;
	timestamp: number;
	transactions: Transaction[];
}

export interface GetChainResponse {
	chain: Block[];
	chain_length: number;
	transaction_pool: Transaction[];
	blockchain_address: string;
	host: string;
	port: number;
	mining: boolean;
	mining_difficulty: number;
	mining_reward: number;
	neighbors: string[];
	wallets: Record<string, Wallet>;
}

export interface GetWalletAmountResponse {
	amount: number;
}

export type ListTransactionPoolsResponse = Transaction[];

export interface Node {
	id: string;
	name: string;
	url: string; // e.g. http://localhost
	port: number;
}

// Optionally, extend Wallet with balance for UI
export interface WalletWithBalance extends Wallet {
	balance?: number;
}
