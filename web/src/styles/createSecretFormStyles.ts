import { makeStyles, Theme } from "@material-ui/core";

export const useStyles = makeStyles((theme: Theme) => ({
	form: {
		display: "flex",
		flexDirection: "column",
		gap: theme.spacing(2),
		maxWidth: 400,
		margin: "auto",
		"& .MuiTextField-root": {
			width: "100%"
		}
	},
	height__full: {
		height: "100%"
	},
	hide: {
		display: "none"
	}
}));
