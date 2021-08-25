import React from "react";
import { Paper, Theme, Typography } from "@material-ui/core";
import MessageIcon from "@material-ui/icons/Message";
import AttachFileIcon from "@material-ui/icons/AttachFile";
import { makeStyles } from "@material-ui/styles";
import Layout from "components/Layout";
import { useHistory } from "react-router";

const useStyles = makeStyles((theme: Theme) => ({
	root: {
		display: "flex",
		justifyContent: "center",
		alignItems: "center",
		width: "100%,",
		height: "100%",
		background: "#ecf0f1"
	},
	container: {
		display: "flex",
		gap: theme.spacing(4),
		flexDirection: "column",
		justifyContent: "center",
		textAlign: "center",
		alignItems: "center",
		maxWidth: "500px",
		width: "100%",
		maxHeight: "80%",
		height: "80%",
		boxShadow: "rgba(0, 0, 0, 0.1) 0px 4px 12px",
		background: "#f5f7f8",
		borderRadius: "5px"
	},
	paper: {
		width: "100%",
		maxWidth: "200px",
		flexWrap: "wrap",
		height: "150px",
		padding: theme.spacing(2),
		cursor: "pointer",
		boxShadow: "rgba(0, 0, 0, 0.1) 0px 4px 12px",
		"&:hover": {
			background: "#f9f9f9",
			transition: "200ms ease"
		}
	},
	card__container: {
		display: "flex",
		gap: theme.spacing(3)
	},
	font__bold: {
		fontWeight: 600
	}
}));

const SecretType: React.FC = () => {
	const classes = useStyles();
	const history = useHistory();

	const handleClick = (type: "file" | "message") => {
		history.push("/create-secret", {
			type
		});
	};

	return (
		<Layout>
			<div className={classes.root}>
				<div className={classes.container}>
					<div>
						<Typography style={{ fontSize: "2rem" }}>
							Welcome on <span style={{ color: "indigo", fontWeight: "bold" }}>WHISPER</span>
						</Typography>
						<Typography variant="body2" className={classes.font__bold}>
							What kind of secret do you want to send ?
						</Typography>
					</div>
					<div className={classes.card__container}>
						<Paper className={classes.paper} elevation={0} onClick={() => handleClick("message")}>
							<div style={{ padding: "1rem 0" }}>
								<MessageIcon fontSize="large" />
							</div>
							<Typography color="primary" className={classes.font__bold}>
								I want to send a message
							</Typography>
						</Paper>
						<Paper className={classes.paper} elevation={0} onClick={() => handleClick("file")}>
							<div style={{ padding: "1rem 0" }}>
								<AttachFileIcon fontSize="large" />
							</div>
							<Typography color="primary" className={classes.font__bold}>
								I want to send a file
							</Typography>
						</Paper>
					</div>
				</div>
			</div>
		</Layout>
	);
};

export default SecretType;
