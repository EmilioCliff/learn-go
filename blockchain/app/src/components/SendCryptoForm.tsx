import { useState } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Alert, AlertDescription } from '@/components/ui/alert';
import { Send, AlertTriangle, CheckCircle } from 'lucide-react';
import { toast } from 'sonner';
import type { Wallet } from '@/lib/types';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { createTransaction } from '@/services/createTransaction';

interface SendCryptoFormProps {
	wallet: Record<string, Wallet>;
	walletAmount: number;
	selectedWalletKey: string;
	selectedNode: string;
	refetchAmount: () => void;
	refetchChain: () => void;
}

export function SendCryptoForm({
	wallet,
	walletAmount,
	selectedWalletKey,
	selectedNode,
	refetchAmount,
	refetchChain,
}: SendCryptoFormProps) {
	const [recipient, setRecipient] = useState('');
	const [amount, setAmount] = useState('');
	const selectedWallet = wallet[selectedWalletKey];

	const [error, setError] = useState('');
	const [success, setSuccess] = useState('');

	const queryClient = useQueryClient();

	const validateForm = () => {
		if (!recipient.trim()) {
			return 'Recipient address is required';
		}
		if (recipient.length < 26 || recipient.length > 35) {
			return 'Invalid address format';
		}
		if (!amount || parseFloat(amount) <= 0) {
			return 'Amount must be greater than 0';
		}
		if (parseFloat(amount) > walletAmount) {
			return 'Insufficient balance';
		}
		return null;
	};

	const mutation = useMutation({
		mutationFn: createTransaction,
		onSuccess: async () => {
			await queryClient.invalidateQueries({
				queryKey: ['chain', 'walletAmount'],
			});
			refetchChain();
			refetchAmount();
			setSuccess(
				`Successfully sent ${amount} BTC to ${recipient.slice(
					0,
					10,
				)}...`,
			);
			setRecipient('');
			setAmount('');

			toast.success('Transaction sent', {
				description: `${amount} BTC sent successfully`,
			});
			setError('');
			setSuccess('');
		},
		onError: (error: any) => {
			toast.error(error.message);
		},
	});

	const handleSubmit = async (e: React.FormEvent) => {
		e.preventDefault();

		const validationError = validateForm();
		if (validationError) {
			setError(validationError);
			return;
		}

		mutation.mutate({
			transaction: {
				sender_private_key: selectedWallet.private_key,
				sender_public_key: selectedWallet.public_key,
				sender_blockchain_address: selectedWallet.blockchain_address,
				recipient_blockchain_address: recipient.trim(),
				value: parseFloat(amount),
			},
			baseUrl: selectedNode,
		});
	};

	return (
		<Card className="elevated-card">
			<CardHeader>
				<CardTitle className="flex items-center gap-2">
					<Send className="h-5 w-5" />
					Send Cryptocurrency
				</CardTitle>
			</CardHeader>
			<CardContent>
				<form onSubmit={handleSubmit} className="space-y-4">
					{/* Recipient Address */}
					<div className="space-y-2">
						<Label htmlFor="recipient">Recipient Address</Label>
						<Input
							id="recipient"
							placeholder="Enter blockchain address (e.g., 1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa)"
							value={recipient}
							onChange={(e) => setRecipient(e.target.value)}
							className="font-mono text-sm"
						/>
					</div>

					{/* Amount */}
					<div className="space-y-2">
						<Label htmlFor="amount">Amount (BTC)</Label>
						<Input
							id="amount"
							type="number"
							step="0.000001"
							min="0"
							// max={walletAmount}
							placeholder="0.000000"
							value={amount}
							onChange={(e) => setAmount(e.target.value)}
						/>
						<div className="text-xs text-muted-foreground">
							Available: {walletAmount.toFixed(6)} BTC
						</div>
					</div>

					{/* Error/Success Messages */}
					{error && (
						<Alert variant="destructive">
							<AlertTriangle className="h-4 w-4" />
							<AlertDescription>{error}</AlertDescription>
						</Alert>
					)}

					{success && (
						<Alert className="border-success bg-success/5">
							<CheckCircle className="h-4 w-4 text-success" />
							<AlertDescription className="text-success">
								{success}
							</AlertDescription>
						</Alert>
					)}

					{/* Send Button */}
					<Button
						type="submit"
						className="w-full"
						disabled={mutation.isPending || !recipient || !amount}
					>
						{mutation.isPending ? (
							<>
								<div className="animate-spin rounded-full h-4 w-4 border-b-2 border-current mr-2" />
								Sending...
							</>
						) : (
							<>
								<Send className="h-4 w-4 mr-2" />
								Send Transaction
							</>
						)}
					</Button>

					{/* Transaction Fee Info */}
					<div className="text-xs text-muted-foreground bg-muted p-3 rounded">
						<div className="flex justify-between">
							<span>Network Fee:</span>
							<span>0.0001 BTC</span>
						</div>
						<div className="flex justify-between">
							<span>Total:</span>
							<span>
								{amount
									? (parseFloat(amount) + 0.0001).toFixed(6)
									: '0.000000'}{' '}
								BTC
							</span>
						</div>
					</div>
				</form>
			</CardContent>
		</Card>
	);
}
