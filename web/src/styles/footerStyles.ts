import { colors, makeStyles } from "@material-ui/core";

export const useStyles = makeStyles({
	root: {
		background: colors["blueGrey"]["900"],
		boxShadow: "0 3px 5px 2px rgba(255, 105, 135, .3)",
		color: "white",
		display: "flex",
		justifyContent: "space-around",
		padding: "1rem 0",
		position: "absolute",
		width: "100%",
		bottom: 0,
		left: 0
	},
	text__white: {
		color: "#fff"
	}
});
