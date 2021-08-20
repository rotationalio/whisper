export interface Secret {
	secret: string;
	password?: string;
	accesses: number;
	lifetime: string;
	filename?: string;
	is_base64: boolean;
	destroyed?: boolean;
	created?: Date;
}
