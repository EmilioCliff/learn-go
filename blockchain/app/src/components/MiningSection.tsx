import { useEffect, useState } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Alert, AlertDescription } from '@/components/ui/alert';
import { Pickaxe, Zap, CheckCircle, Clock, Coins, Square } from 'lucide-react';
import { toast } from 'sonner';
import type { Wallet } from '@/lib/types';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { getWalletAmount } from '@/services/getWallentAmount';
import { mine } from '@/services/mine';
import { startMining } from '@/services/startMining';
import { stopMining } from '@/services/stopMining';

interface MiningSectionProps {
	minerWallet: Wallet;
	mining: boolean;
	mining_difficulty: number;
	mining_reward: number;
	selectedNode: string;
	refetchChain: () => void;
}

export function MiningSection({
	minerWallet,
	mining,
	mining_difficulty,
	mining_reward,
	selectedNode,
	refetchChain,
}: MiningSectionProps) {
	const [lastBlockReward, setLastBlockReward] = useState<number | null>(null);

	const queryClient = useQueryClient();

	useEffect(() => {
		setLastBlockReward(null);
	}, [selectedNode]);

	const { data: walletAmountData, refetch: refetchAmount } = useQuery({
		queryKey: [
			'walletAmount',
			minerWallet.blockchain_address,
			selectedNode,
		],
		queryFn: () =>
			getWalletAmount(minerWallet.blockchain_address, selectedNode),
		refetchInterval: mining ? 5000 : false,
	});

	const mineMutation = useMutation({
		mutationFn: mine,
		onSuccess: async () => {
			await queryClient.invalidateQueries({
				queryKey: ['walletAmount'],
			});
			refetchAmount();
			refetchChain();
			toast.success('Block mined successfully!', {
				description: `Earned ${mining_reward} BTC`,
			});
			setLastBlockReward(mining_reward);
		},
		onError: (error: any) => {
			toast.error(error.message);
		},
	});

	const startMiningMutation = useMutation({
		mutationFn: startMining,
		onSuccess: async () => {
			await queryClient.invalidateQueries({
				queryKey: ['walletAmount'],
			});
			setLastBlockReward(null);
			toast.success('Mining will continue!', {});
			refetchChain();
		},
		onError: (error: any) => {
			toast.error(error.message);
		},
	});

	const stopMiningMutation = useMutation({
		mutationFn: stopMining,
		onSuccess: async () => {
			await queryClient.invalidateQueries({
				queryKey: ['walletAmount'],
			});
			setLastBlockReward(null);
			toast.success('Mining has been stoped!', {});
			refetchChain();
		},
		onError: (error: any) => {
			toast.error(error.message);
		},
	});

	return (
		<Card className="elevated-card">
			<CardHeader>
				<CardTitle className="flex items-center gap-2">
					<Pickaxe className="h-5 w-5" />
					Mining
					{mining ? (
						<Button
							onClick={() =>
								stopMiningMutation.mutate(selectedNode)
							}
							className="ml-auto"
							size="sm"
							variant="destructive"
						>
							<Square className="h-5 w-5" />
							Stop Mining
						</Button>
					) : (
						<Button
							onClick={() =>
								startMiningMutation.mutate(selectedNode)
							}
							className="ml-auto"
							size="sm"
							variant="outline"
						>
							<Pickaxe className="h-5 w-5" />
							Start Mining
						</Button>
					)}
				</CardTitle>
			</CardHeader>
			<CardContent className="space-y-4">
				{/* Miner Wallet Info */}
				<div className="bg-gradient-to-r from-crypto/10 to-crypto-light/10 p-4 rounded-lg border">
					<div className="text-sm text-muted-foreground">
						Miner Balance
					</div>
					<div className="text-lg font-bold flex items-center gap-2">
						<Coins className="h-4 w-4" />
						{walletAmountData?.amount.toFixed(6)} BTC
					</div>
					<div className="text-xs text-muted-foreground mt-1 crypto-address">
						{minerWallet.blockchain_address}
					</div>
				</div>

				{/* Mining Status */}
				<div className="space-y-3">
					{mineMutation.isPending ? (
						<div className="space-y-2">
							<div className="flex items-center justify-between">
								<span className="text-sm font-medium">
									Mining in progress...
								</span>
								<Badge variant="default" className="gap-1">
									<Zap className="h-3 w-3" />
									Active
								</Badge>
							</div>
						</div>
					) : (
						<div className="text-center space-y-2">
							<Badge
								variant={mining ? 'default' : 'secondary'}
								className="gap-1"
							>
								<Clock className="h-3 w-3" />
								{mining ? 'Auto Mining' : 'Not mining'}
							</Badge>
							{!mining && (
								<div className="text-sm text-muted-foreground">
									Ready to mine the next block
								</div>
							)}
						</div>
					)}
				</div>

				{/* Last Block Reward */}
				{lastBlockReward && (
					<Alert className="border-success bg-success/5">
						<CheckCircle className="h-4 w-4 text-success" />
						<AlertDescription className="text-success">
							Last block reward: {lastBlockReward.toFixed(4)} BTC
						</AlertDescription>
					</Alert>
				)}

				{/* Mining Controls */}
				{!mining && !mineMutation.isPending && (
					<div className="space-y-2">
						<Button
							className="w-full bg-gradient-to-r from-crypto to-crypto"
							onClick={() => mineMutation.mutate(selectedNode)}
						>
							<Pickaxe className="h-4 w-4 mr-2" />
							Start Mining
						</Button>
					</div>
				)}
				{!mining && mineMutation.isPending && (
					<div className="space-y-2">
						<Button
							variant={'destructive'}
							className="w-full "
							onClick={() => mineMutation.reset()}
						>
							<Pickaxe className="h-4 w-4 mr-2" />
							Stop Mining
						</Button>
					</div>
				)}

				{/* Mining Stats */}
				<div className="bg-muted p-3 rounded space-y-2">
					<div className="text-sm font-medium text-muted-foreground">
						Mining Statistics
					</div>
					<div className="grid grid-cols-2 gap-2 text-xs">
						<div>
							<span className="text-muted-foreground">
								Difficulty:
							</span>
							<div className="font-medium">
								{mining_difficulty}
							</div>
						</div>
						<div>
							<span className="text-muted-foreground">
								Reward:
							</span>
							<div className="font-medium">{mining_reward}</div>
						</div>
					</div>
				</div>
			</CardContent>
		</Card>
	);
}
