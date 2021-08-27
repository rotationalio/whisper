import { makeStyles } from "@material-ui/core";

export const useStyles = makeStyles(() => ({
	root: {
		flexGrow: 1
	},
	h__full: {
		height: "100%",
		boxShadow: "rgba(0, 0, 0, 0.1) 0px 4px 12px",
		background: "#f5f7f8",
		borderRadius: "5px",
		flexDirection: "column",
		justifyContent: "center",
		alignItems: "center",
		width: "100%"
	},
	w__full: {
		width: "100%"
	},
	alert: {
		maxWidth: 400,
		margin: "0 auto",
		boxSizing: "border-box",
		marginBottom: "1rem"
	}
}));
