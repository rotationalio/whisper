import { colors, makeStyles } from "@material-ui/core";

export const useStyles = makeStyles({
	root: {
		background: colors["blueGrey"]["900"],
		boxShadow: "0 3px 5px 2px rgba(255, 105, 135, .3)",
		color: "white",
		display: "flex",
		justifyContent: "space-around",
		padding: "1rem 0",
		width: "100%"
	},
	text__white: {
		color: "#fff"
	}
});
