import { useState } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import {
	Select,
	SelectContent,
	SelectItem,
	SelectTrigger,
	SelectValue,
} from '@/components/ui/select';
import { Badge } from '@/components/ui/badge';
import {
	Wallet as WalletIcon,
	Eye,
	EyeOff,
	Copy,
	RefreshCw,
	Plus,
} from 'lucide-react';
import { toast } from 'sonner';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import type { Wallet } from '@/lib/types';
import { createWallet } from '@/services/createWallet';

interface WalletSectionProps {
	wallet: Record<string, Wallet>;
	miners_address: string;
	selectedNode: string;
	selectedWalletKey: string;
	walletAmount: number;
	refetchAmount: () => void;
	setSelectedWalletKey: (key: string) => void;
}

export function WalletSection({
	wallet,
	miners_address,
	selectedNode,
	selectedWalletKey,
	walletAmount,
	setSelectedWalletKey,
	refetchAmount,
}: WalletSectionProps) {
	const [showPrivateKey, setShowPrivateKey] = useState(false);
	const selectedWallet = wallet[selectedWalletKey];

	const queryClient = useQueryClient();

	// useEffect(() => {
	// 	setSelectedWalletKey(miners_address);
	// }, [selectedNode, miners_address]);

	const mutation = useMutation({
		mutationFn: createWallet,
		onSuccess: async (data) => {
			await queryClient.invalidateQueries({
				queryKey: ['chain', 'walletAmount'],
			});
			toast.success('Wallet created successfully');
			wallet[data.blockchain_address] = data;
			setSelectedWalletKey(data.blockchain_address);
			refetchAmount();
		},
		onError: (error: any) => {
			toast.error(error.message);
		},
	});

	const onWalletChange = (value: string) => {
		setSelectedWalletKey(value);
		setShowPrivateKey(false);
	};

	const copyToClipboard = (text: string, label: string) => {
		navigator.clipboard.writeText(text);
		toast('Copied!', {
			description: `${label} copied to clipboard`,
		});
	};

	return (
		<Card className="elevated-card">
			<CardHeader>
				<CardTitle className="flex items-center gap-2">
					<WalletIcon className="h-5 w-5" />
					Wallet
					<Button
						onClick={() => mutation.mutate(selectedNode)}
						className="ml-auto"
						size="sm"
						variant="outline"
					>
						<Plus className="h-5 w-5" />
						Create Wallet
					</Button>
				</CardTitle>
			</CardHeader>
			<CardContent className="space-y-4">
				{/* Wallet Selection */}
				<div>
					<label className="text-sm font-medium mb-2 block">
						Active Wallet
					</label>
					<Select
						value={selectedWalletKey}
						onValueChange={onWalletChange}
					>
						<SelectTrigger>
							<SelectValue />
						</SelectTrigger>
						<SelectContent>
							{Object.entries(wallet).map(([key, wallet]) => (
								<SelectItem key={key} value={key}>
									{key === miners_address
										? 'Miner Wallet'
										: 'User Wallet'}{' '}
									({wallet.blockchain_address.slice(0, 6)}...
									{wallet.blockchain_address.slice(-4)})
								</SelectItem>
							))}
						</SelectContent>
					</Select>
				</div>

				{/* Balance */}
				<div className="bg-gradient-to-r from-crypto/10 to-crypto-light/10 p-4 rounded-lg border">
					<div className="text-sm text-muted-foreground">Balance</div>
					<div className="text-2xl font-bold flex items-center gap-2">
						{walletAmount.toFixed(6)} BTC
						<Button
							variant="ghost"
							size="sm"
							onClick={() => {
								refetchAmount();
								toast('Balance refreshed');
							}}
						>
							<RefreshCw className="h-4 w-4" />
						</Button>
					</div>
				</div>

				{/* Wallet Details */}
				<div className="space-y-3">
					{/* Public Key */}
					<div>
						<label className="text-sm font-medium text-muted-foreground">
							Public Key
						</label>
						<div className="flex items-center gap-2 mt-1">
							<div className="crypto-address flex-1">
								{selectedWallet?.public_key}
							</div>
							<Button
								variant="ghost"
								size="sm"
								onClick={() =>
									copyToClipboard(
										selectedWallet?.public_key,
										'Public key',
									)
								}
							>
								<Copy className="h-4 w-4" />
							</Button>
						</div>
					</div>

					{/* Private Key */}
					<div>
						<label className="text-sm font-medium text-muted-foreground">
							Private Key
						</label>
						<div className="flex items-center gap-2 mt-1">
							<div className="crypto-address flex-1">
								{showPrivateKey
									? selectedWallet?.private_key
									: '••••••••••••••••••••••••••••••••••••••••••••••••••••'}
							</div>
							<Button
								variant="ghost"
								size="sm"
								onClick={() =>
									setShowPrivateKey(!showPrivateKey)
								}
							>
								{showPrivateKey ? (
									<EyeOff className="h-4 w-4" />
								) : (
									<Eye className="h-4 w-4" />
								)}
							</Button>
							{showPrivateKey && (
								<Button
									variant="ghost"
									size="sm"
									onClick={() =>
										copyToClipboard(
											selectedWallet?.private_key,
											'Private key',
										)
									}
								>
									<Copy className="h-4 w-4" />
								</Button>
							)}
						</div>
					</div>

					{/* Address */}
					<div>
						<label className="text-sm font-medium text-muted-foreground">
							Blockchain Address
						</label>
						<div className="flex items-center gap-2 mt-1">
							<div className="crypto-address flex-1">
								{selectedWallet?.blockchain_address}
							</div>
							<Button
								variant="ghost"
								size="sm"
								onClick={() =>
									copyToClipboard(
										selectedWallet?.blockchain_address,
										'Address',
									)
								}
							>
								<Copy className="h-4 w-4" />
							</Button>
						</div>
					</div>
				</div>

				{/* Wallet Type Badge */}
				<div className="pt-2">
					<Badge
						variant={
							selectedWalletKey === miners_address
								? 'default'
								: 'secondary'
						}
					>
						{selectedWalletKey === miners_address
							? 'Mining Wallet'
							: 'User Wallet'}
					</Badge>
				</div>
			</CardContent>
		</Card>
	);
}
