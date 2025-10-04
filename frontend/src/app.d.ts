declare global {
	namespace App {
		interface Locals {
			user: { id: number } | null;
			token: string | null;
			activeBandId: string | undefined;
		}
	}
	interface RequestInit {
		duplex?: 'half';
	}
}

export {};