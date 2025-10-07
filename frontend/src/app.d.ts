declare global {
	namespace App {
		interface Locals {
			user: {
				id: number;
				username: string;
				band_name: string;
				role: string;
			} | null;
			token: string | null;
			activeBandId: string | undefined;
		}
	}
	interface RequestInit {
		duplex?: 'half';
	}
}

export {};