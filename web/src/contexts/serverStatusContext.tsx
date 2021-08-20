import React from "react";
import { ServerStatus } from "types/ServerStatus";

type Status = {
	status: ServerStatus;
	timestamp: Date;
	version: string;
	host: string;
};

const ServerStatusContext = React.createContext<
	[Partial<Status>, React.Dispatch<React.SetStateAction<Partial<Status>>>] | undefined
>(undefined);

const ServerStatusProvider: React.FC = props => {
	const [status, setStatus] = React.useState<Partial<Status>>({
		status: undefined,
		timestamp: undefined,
		version: undefined,
		host: undefined
	});
	return <ServerStatusContext.Provider value={[status, setStatus]} {...props} />;
};

const useServerStatus = (): [Partial<Status>, React.Dispatch<React.SetStateAction<Partial<Status>>>] => {
	const context = React.useContext(ServerStatusContext);
	if (!context) {
		throw new Error("useServerStatus should be used within a ServerStatusProvider");
	}

	return context;
};

export { ServerStatusProvider, useServerStatus };
