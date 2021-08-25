import { Lifetime } from "utils/interfaces";

const FILE_SIZE = 64 * 1024;

const LIFETIME_OPTIONS: Lifetime[] = [
	{ value: "5m", label: "5 min" },
	{ value: "15m", label: "15 min" },
	{ value: "30m", label: "30 min" },
	{ value: "1h", label: "1 hour" },
	{ value: "2h", label: "2 hours" },
	{ value: "3h", label: "3 hours" },
	{ value: "24h", label: "1 day" },
	{ value: "48h", label: "2 days" },
	{ value: "72h", label: "3 days" },
	{ value: "168h", label: "7 days" }
];

export { FILE_SIZE, LIFETIME_OPTIONS };
