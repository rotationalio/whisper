import { Box, Button, makeStyles, Typography } from "@material-ui/core";
import { Alert, AlertProps, AlertTitle } from "@material-ui/lab";
import { Link } from "react-router-dom";

const useStyles = makeStyles(() => ({
	root: {
		borderRadius: "0 !important",
		"& > .MuiAlert-message > p::first-letter": {
			textTransform: "uppercase"
		}
	},
	container: {
		position: "absolute",
		top: "50%",
		left: "50%",
		transform: "translate(-50%, -50%)",
		width: "100%",
		maxWidth: "500px"
	},
	box: {
		width: "100%",
		border: "1px dashed red",
		margin: "0 auto"
	},
	link: {
		textDecoration: "none"
	}
}));

interface SecretNotFoundProps extends AlertProps {
	title: string;
	message: string;
}

const SecretNotFound: React.FC<SecretNotFoundProps> = ({ message, title, ...props }) => {
	const classes = useStyles();

	return (
		<div className={classes.container}>
			<Box className={classes.box}>
				<Alert severity="error" className={classes.root} {...props}>
					<AlertTitle>{title}</AlertTitle>
					<Typography>{message ? `${message}.` : "No secret exists with the specified token."}</Typography>
				</Alert>
			</Box>
			<div style={{ marginTop: "2rem" }}>
				<Link to="/" className={classes.link}>
					<Button variant="contained" color="primary">
						Create a new secret
					</Button>
				</Link>
			</div>
		</div>
	);
};

export default SecretNotFound;
