import { makeStyles, Theme } from "@material-ui/core";

export const useStyles = makeStyles((theme: Theme) => ({
	root: {
		flexGrow: 1,
		padding: theme.spacing(2)
	},
	h__full: {
		height: "100%"
	},
	w__full: {
		width: "100%",
		position: "absolute",
		top: "50%",
		left: "50%",
		transform: "translate(-50%, -50%)"
	},
	alert: {
		maxWidth: 400,
		margin: "0 auto",
		boxSizing: "border-box",
		marginBottom: "1rem"
	}
}));
