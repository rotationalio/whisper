export type Status = "ok" | "unhealthy" | "maintainance";

export interface ServerStatus {
	status: Status;
	timestamp: Date;
	version: string;
}
