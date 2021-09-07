import React from "react";
import { makeStyles, Theme } from "@material-ui/core";

const useStyles = makeStyles((theme: Theme) => ({
	root: {
		flexGrow: 1
	},
	container: {
		boxShadow: "rgba(0, 0, 0, 0.1) 0px 4px 12px",
		background: "#fff",
		padding: theme.spacing(2),
		borderRadius: "5px"
	},
	chpisContainer: {
		marginBottom: 5,
		textAlign: "end"
	}
}));

type ContentLayoutProps = {
	children: React.ReactNode;
};

const ContentLayout: React.FC<ContentLayoutProps> = ({ children }) => {
	const classes = useStyles();

	return (
		<div className={classes.root}>
			<div className={classes.container}>{children}</div>
		</div>
	);
};

export default ContentLayout;
