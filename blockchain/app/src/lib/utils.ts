import { clsx, type ClassValue } from 'clsx';
import { format } from 'date-fns';
import { twMerge } from 'tailwind-merge';

export function cn(...inputs: ClassValue[]) {
	return twMerge(clsx(inputs));
}

export function formatTimestamp(timestamp: number): string {
	// if (!timestamp) return 'N/A';
	let ts = timestamp;

	// normalize to ms
	if (ts > 1e15) {
		ts = Math.floor(ts / 1e6); // nanoseconds → ms
	} else if (ts < 1e12) {
		ts = ts * 1000; // seconds → ms
	}

	const date = new Date(ts);
	if (isNaN(date.getTime())) return 'Invalid date';

	// Example: 9/24/2025, 9:14:43 PM
	return format(date, 'M/d/yyyy, h:mm:ss a');
}

export function stripProtocol(url: string): string {
	if (url.startsWith('http://')) {
		return url.slice(7);
	}
	if (url.startsWith('https://')) {
		return url.slice(8);
	}
	return url;
}
