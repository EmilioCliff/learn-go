import { useState } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { ScrollArea } from '@/components/ui/scroll-area';
import { Separator } from '@/components/ui/separator';
import { Clock, RefreshCw, ArrowRight, AlertCircle } from 'lucide-react';
import type { Transaction } from '@/lib/types';

interface PendingTransactionsProps {
	transactions: Transaction[];
	refreshChain: () => void;
}

export function PendingTransactions({
	transactions,
	refreshChain,
}: PendingTransactionsProps) {
	const [isRefreshing, setIsRefreshing] = useState(false);

	const refreshPendingTransactions = async () => {
		setIsRefreshing(true);
		await new Promise((resolve) => setTimeout(resolve, 1000));
		refreshChain();
		setIsRefreshing(false);
	};

	const formatTimestamp = (timestamp: number) => {
		const diff = Date.now() - timestamp;
		const minutes = Math.floor(diff / 60000);
		const seconds = Math.floor((diff % 60000) / 1000);

		if (minutes > 0) {
			return `${minutes}m ${seconds}s ago`;
		}
		return `${seconds}s ago`;
	};

	const formatAddress = (address: string) => {
		return `${address.slice(0, 8)}...${address.slice(-8)}`;
	};

	const getTotalValue = () => {
		return transactions.reduce((sum, tx) => sum + tx.value, 0);
	};

	return (
		<Card className="elevated-card">
			<CardHeader>
				<div className="flex items-center justify-between">
					<CardTitle className="flex items-center gap-2">
						<Clock className="h-5 w-5" />
						Pending Transactions ({transactions.length})
					</CardTitle>
					<Button
						variant="outline"
						size="sm"
						onClick={refreshPendingTransactions}
						disabled={isRefreshing}
					>
						<RefreshCw
							className={`h-4 w-4 ${
								isRefreshing ? 'animate-spin' : ''
							}`}
						/>
					</Button>
				</div>
			</CardHeader>
			<CardContent>
				{/* Summary */}
				<div className="bg-gradient-to-r from-warning/10 to-warning/5 p-4 rounded-lg border mb-4">
					<div className="flex items-center justify-between">
						<div>
							<div className="text-sm text-muted-foreground">
								Total Pending Value
							</div>
							<div className="text-lg font-bold">
								{getTotalValue().toFixed(6)} BTC
							</div>
						</div>
						<Badge variant="outline" className="gap-1">
							<AlertCircle className="h-3 w-3" />
							TransactionPool
						</Badge>
					</div>
				</div>

				{/* Transactions List */}
				{transactions.length === 0 ? (
					<div className="text-center py-8 text-muted-foreground">
						<Clock className="h-8 w-8 mx-auto mb-2 opacity-50" />
						<div>No pending transactions</div>
						<div className="text-xs">
							All transactions have been confirmed
						</div>
					</div>
				) : (
					<ScrollArea className="h-[300px]">
						<div className="space-y-3">
							{transactions.map((tx, index) => (
								<div key={index}>
									<div className="bg-background/50 border rounded-lg p-3">
										{/* Transaction Header */}
										<div className="flex items-center justify-between mb-2">
											<Badge
												variant="outline"
												className="text-xs"
											>
												{formatTimestamp(
													tx.timestamp || 0,
												)}
											</Badge>
											<Badge
												variant="secondary"
												className="text-xs"
											>
												{tx.value} BTC
											</Badge>
										</div>

										{/* Transaction Flow */}
										<div className="flex items-center gap-2 text-sm">
											<div className="crypto-address text-xs flex-1">
												{formatAddress(
													tx.sender_blockchain_address,
												)}
											</div>
											<ArrowRight className="h-3 w-3 text-muted-foreground flex-shrink-0" />
											<div className="crypto-address text-xs flex-1">
												{formatAddress(
													tx.recipient_blockchain_address,
												)}
											</div>
										</div>

										{/* Transaction ID */}
										<div className="mt-2 text-xs text-muted-foreground">
											ID:{' '}
											<span className="crypto-hash">
												pending_tx_{index + 1}
											</span>
										</div>
									</div>

									{index < transactions.length - 1 && (
										<Separator className="my-2" />
									)}
								</div>
							))}
						</div>
					</ScrollArea>
				)}

				{/* Info */}
				<div className="mt-4 text-xs text-muted-foreground bg-muted p-3 rounded">
					<div className="flex items-center gap-1 mb-1">
						<AlertCircle className="h-3 w-3" />
						<span className="font-medium">
							TransactionPool Info
						</span>
					</div>
					<div>
						Transactions are waiting to be included in the next
						block. Mining a new block will confirm some of these
						transactions.
					</div>
				</div>
			</CardContent>
		</Card>
	);
}
