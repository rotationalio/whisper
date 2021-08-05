export type LifetimeValue = "5m" | "15m" | "30m" | "1h" | "2h" | "3h" | "24h" | "48h" | "72h" | "168h";

export interface Lifetime {
	label: string;
	value: LifetimeValue;
}
