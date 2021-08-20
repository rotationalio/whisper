import { Box, Link, Tooltip, Typography } from "@material-ui/core";
import { useServerStatus } from "contexts/serverStatusContext";
import { useEffect, useState } from "react";
import { Status } from "utils/enums/status";
import * as footerStyles from "../styles/footerStyles";
import Badge from "./Badge";

type StatusColor = "green" | "yellow" | "red";

const Footer: React.FC = () => {
	const [statusColor, setStatusColor] = useState<StatusColor>("green");
	const classes = footerStyles.useStyles();
	const [status] = useServerStatus();

	useEffect(() => {
		switch (status?.status) {
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
	}, [status]);

	return (
		<footer className={classes.root}>
			<Typography>
				Made with &spades; by{" "}
				<Link href="https://rotational.io" target="_blank" className={classes.text__white}>
					Rotational Labs
				</Link>
			</Typography>

			<Box display="flex" alignItems="center" gridGap=".5rem">
				<Tooltip title={`connected to ${status.host}`} aria-label="add" style={{ cursor: "pointer" }}>
					<Box display="flex" gridGap="1rem">
						<Badge color={statusColor} content="status" />
						<Typography variant="caption">version: {status?.version ? status?.version : "0.0.0"}</Typography>
					</Box>
				</Tooltip>
			</Box>
		</footer>
	);
};

export default Footer;
