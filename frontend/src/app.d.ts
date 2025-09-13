declare global {
	namespace App {
		interface Locals {
			user: { id: number; bandId: number } | null;
			token: string | null;
		}
	}
	interface RequestInit {
		duplex?: 'half';
	}
}

export {};