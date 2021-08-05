export interface Secret {
	secret: string;
	password: string;
	accesses: number;
	lifetime: string;
	filename?: string;
	is_base64: false;
}
