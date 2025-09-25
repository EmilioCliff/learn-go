import { useState, useEffect } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import {
	Select,
	SelectContent,
	SelectItem,
	SelectTrigger,
	SelectValue,
} from '@/components/ui/select';
import { Separator } from '@/components/ui/separator';
import { WalletSection } from './WalletSection';
import { SendCryptoForm } from './SendCryptoForm';
import { BlockchainView } from './BlockchainView';
import { MiningSection } from './MiningSection';
import { PendingTransactions } from './PendingTransactions';
import { toast } from 'sonner';
import { Server } from 'lucide-react';
import {
	keepPreviousData,
	useQuery,
	useQueryClient,
} from '@tanstack/react-query';
import { getChain } from '@/services/getChain';
import { DEFAULT_BASE_URL } from '@/lib/constants';
import { stripProtocol } from '@/lib/utils';
import { getWalletAmount } from '@/services/getWallentAmount';

export function BlockchainDashboard() {
	const [selectedNode, setSelectedNode] = useState<string>(DEFAULT_BASE_URL);
	const [selectedWalletKey, setSelectedWalletKey] = useState('');

	const queryClient = useQueryClient();

	const { data: walletAmountData, refetch: refetchAmount } = useQuery({
		queryKey: ['walletAmount', selectedWalletKey, selectedNode],
		queryFn: () => getWalletAmount(selectedWalletKey, selectedNode),
	});

	const {
		data: chainData,
		isLoading: isChainLoading,
		refetch: refetchChain,
	} = useQuery({
		queryKey: ['chain', selectedNode],
		queryFn: () => getChain(selectedNode),
		placeholderData: keepPreviousData,
		// refetchInterval: 5000,
	});

	useEffect(() => {
		if (chainData && chainData.blockchain_address) {
			setSelectedWalletKey(chainData.blockchain_address);
		}
	}, [chainData]);

	const handleNodeChange = (node: string) => {
		if (node) {
			setSelectedNode(node);
			toast('Node switched', {
				description: `Connected to ${node})`,
			});
			queryClient.invalidateQueries();
		}
	};

	if (isChainLoading) return <div>Loading...</div>;

	return (
		<div className="min-h-screen bg-background p-6">
			<div className="max-w-7xl mx-auto space-y-6">
				{/* Header */}
				<div className="flex items-center justify-between">
					<div>
						<h1 className="text-3xl font-bold tracking-tight">
							Blockchain Demo
						</h1>
						<p className="text-muted-foreground">
							Multi-node blockchain network interface
						</p>
					</div>
				</div>

				{/* Node Selection */}
				<Card className="gradient-card">
					<CardHeader>
						<CardTitle className="flex items-center gap-2">
							<Server className="h-5 w-5" />
							Network Node
						</CardTitle>
					</CardHeader>
					<CardContent>
						<div className="flex items-center gap-4">
							<Select
								value={selectedNode}
								onValueChange={handleNodeChange}
							>
								<SelectTrigger className="w-[300px]">
									<SelectValue />
								</SelectTrigger>
								<SelectContent>
									{chainData?.neighbors.map((node, idx) => (
										<SelectItem key={idx} value={node}>
											{/* {node.name} ({node.url}:{node.port}) */}
											{stripProtocol(node)}
										</SelectItem>
									))}
								</SelectContent>
							</Select>
							<div className="text-sm text-muted-foreground">
								Connected to:{' '}
								<span className="crypto-mono">
									{selectedNode}
								</span>
							</div>
						</div>
					</CardContent>
				</Card>

				<div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
					{/* Wallet Section */}
					<div className="lg:col-span-1">
						<WalletSection
							wallet={chainData?.wallets || {}}
							miners_address={chainData?.blockchain_address || ''}
							selectedNode={selectedNode}
							selectedWalletKey={selectedWalletKey}
							walletAmount={walletAmountData?.amount || 0}
							refetchAmount={refetchAmount}
							setSelectedWalletKey={setSelectedWalletKey}
						/>
					</div>

					{/* Send Crypto */}
					<div className="lg:col-span-1">
						<SendCryptoForm
							wallet={chainData?.wallets || {}}
							walletAmount={walletAmountData?.amount || 0}
							selectedWalletKey={selectedWalletKey}
							refetchAmount={refetchAmount}
							selectedNode={selectedNode}
							refetchChain={refetchChain}
						/>
					</div>

					{/* Mining Section */}
					<div className="lg:col-span-1">
						<MiningSection
							minerWallet={
								chainData?.wallets[
									chainData.blockchain_address
								] || {
									private_key: '',
									public_key: '',
									blockchain_address: '',
								}
							}
							mining={chainData?.mining || false}
							mining_difficulty={
								chainData?.mining_difficulty || 0
							}
							refetchChain={refetchChain}
							mining_reward={chainData?.mining_reward || 0}
							selectedNode={selectedNode}
						/>
					</div>
				</div>

				<Separator />

				{/* Blockchain and Transactions */}
				<div className="grid gap-6 lg:grid-cols-2">
					{chainData && (
						<BlockchainView
							refreshChain={refetchChain}
							data={chainData}
						/>
					)}
					<PendingTransactions
						transactions={chainData?.transaction_pool || []}
						refreshChain={refetchChain}
					/>
				</div>
			</div>
		</div>
	);
}
