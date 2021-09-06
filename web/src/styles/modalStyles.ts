import { createStyles, makeStyles, Theme } from "@material-ui/core";

export const useStyles = makeStyles((theme: Theme) =>
	createStyles({
		paper: {
			position: "absolute",
			minWidth: 500,
			backgroundColor: theme.palette.background.paper,
			top: "50%",
			left: "50%",
			transform: "translate(-50%, -50%)",
			padding: theme.spacing(2),
			display: "flex",
			flexDirection: "column",
			gap: "2rem"
		},
		ellipsis: {
			whiteSpace: "nowrap",
			textOverflow: "ellipsis",
			display: "block",
			overflow: "scroll",
			scrollbarWidth: "none",
			msOverflowStyle: "none",
			"&::-webkit-scrollbar": {
				display: "none"
			}
		},
		close: {
			position: "absolute",
			right: 15
		}
	})
);
