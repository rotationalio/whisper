import { Box, Link, Tooltip, Typography } from "@material-ui/core";
import { AxiosError } from "axios";
import { useEffect, useState } from "react";
import getStatus from "services/status";
import { Status } from "utils/enums/status";
import { ServerStatus } from "utils/interfaces/Status";
import * as footerStyles from "../styles/footerStyles";
import Badge from "./Badge";

type StatusColor = "green" | "yellow" | "red";

const Footer: React.FC = () => {
	const [serverStatus, setServerStatus] = useState<ServerStatus>();
	const [statusColor, setStatusColor] = useState<StatusColor>("green");
	const [endpoint, setEndpoint] = useState<string>("");
	const classes = footerStyles.useStyles();

	useEffect(() => {
		getStatus()
			.then(response => {
				setServerStatus(response.data);

				switch (response.data.status) {
					case Status.ok:
						setStatusColor("green");
						break;
					case Status.unhealthy:
						setStatusColor("red");
						break;
					case Status.maintainance:
						setStatusColor("yellow");
						break;
				}

				if (typeof response.config.url === "string") {
					setEndpoint(response?.config?.url);
				}
			})
			.catch((error: AxiosError) => {
				console.error("[Footer] fetch server status", error.message);
			});
	}, []);

	return (
		<footer className={classes.root}>
			<Typography>
				Made with &spades; by{" "}
				<Link href="https://rotational.io" target="_blank" className={classes.text__white}>
					Rotational Labs
				</Link>
			</Typography>

			{/* <Typography>Made with &spades; by Rotational Labs</Typography> */}
			<Box display="flex" alignItems="center" gridGap=".5rem">
				<Tooltip title={`connected to ${endpoint}`} aria-label="add" style={{ cursor: "pointer" }}>
					<Box display="flex" gridGap="1rem">
						<Badge color={statusColor} content="status" />
						<Typography variant="caption">version: {serverStatus?.version}</Typography>
					</Box>
				</Tooltip>
			</Box>
		</footer>
	);
};

export default Footer;
