import { useState } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { ScrollArea } from '@/components/ui/scroll-area';
import { Blocks, RefreshCw, ArrowRight, Clock, Hash } from 'lucide-react';
import { formatTimestamp } from '@/lib/utils';
import type { GetChainResponse } from '@/lib/types';

export function BlockchainView({
	data,
	refreshChain,
}: {
	data: GetChainResponse;
	refreshChain: () => void;
}) {
	const [isRefreshing, setIsRefreshing] = useState(false);

	const refreshBlockchain = async () => {
		setIsRefreshing(true);
		await new Promise((resolve) => setTimeout(resolve, 1000));
		refreshChain();
		setIsRefreshing(false);
	};

	const formatHash = (hash: string) => {
		return `${hash.slice(0, 10)}...${hash.slice(-10)}`;
	};

	return (
		<Card className="elevated-card">
			<CardHeader>
				<div className="flex items-center justify-between">
					<CardTitle className="flex items-center gap-2">
						<Blocks className="h-5 w-5" />
						Blockchain ({data.chain_length} blocks)
					</CardTitle>
					<Button
						variant="outline"
						size="sm"
						onClick={refreshBlockchain}
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
				<ScrollArea className="h-[400px]">
					<div className="space-y-4">
						{data.chain.map((block, index) => (
							<div key={index} className="relative">
								{/* Block Card */}
								<div className="bg-gradient-to-r from-card to-muted/50 border rounded-lg p-4">
									{/* Block Header */}
									<div className="flex items-center justify-between mb-3">
										<div className="flex items-center gap-2">
											<Badge
												variant={
													index === 0
														? 'default'
														: 'secondary'
												}
											>
												Block #{index}
											</Badge>
											{index === 0 && (
												<Badge variant="outline">
													Genesis
												</Badge>
											)}
										</div>
										<div className="text-xs text-muted-foreground flex items-center gap-1">
											<Clock className="h-3 w-3" />
											{formatTimestamp(block.timestamp)}
										</div>
									</div>

									{/* Block Details */}
									<div className="grid grid-cols-1 gap-2 text-sm">
										<div className="flex items-center gap-2">
											<Hash className="h-3 w-3 text-muted-foreground" />
											<span className="text-muted-foreground">
												Previous:
											</span>
											<span className="crypto-hash font-mono">
												{formatHash(
													block.previous_hash,
												)}
											</span>
										</div>

										<div className="flex items-center gap-4">
											<span className="text-muted-foreground">
												Nonce:{' '}
												<span className="crypto-mono">
													{block.nonce}
												</span>
											</span>
											<span className="text-muted-foreground">
												Transactions:{' '}
												<span className="font-medium">
													{block.transactions.length}
												</span>
											</span>
										</div>
									</div>

									{/* Transactions */}
									{block.transactions.length > 0 && (
										<div className="mt-3 pt-3 border-t">
											<div className="text-xs font-medium text-muted-foreground mb-2">
												Transactions:
											</div>
											<div className="space-y-2">
												{block.transactions.map(
													(tx, idx) => (
														<div
															key={idx}
															className="bg-background/50 rounded p-2 text-xs"
														>
															<div className="flex items-center justify-between">
																<span className="crypto-address">
																	{formatHash(
																		tx.sender_blockchain_address,
																	)}
																</span>
																<ArrowRight className="h-3 w-3 text-muted-foreground" />
																<span className="crypto-address">
																	{formatHash(
																		tx.recipient_blockchain_address,
																	)}
																</span>
																<Badge
																	variant="outline"
																	className="text-xs"
																>
																	{tx.value}{' '}
																	BTC
																</Badge>
															</div>
														</div>
													),
												)}
											</div>
										</div>
									)}
								</div>

								{/* Connection Line */}
								{index < data.chain.length - 1 && (
									<div className="absolute left-1/2 -bottom-2 w-0.5 h-4 bg-border transform -translate-x-1/2" />
								)}
							</div>
						))}
					</div>
				</ScrollArea>
			</CardContent>
		</Card>
	);
}
