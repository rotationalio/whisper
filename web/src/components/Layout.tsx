import { makeStyles, Theme } from "@material-ui/core";

type LayoutProps = {
	children: React.ReactNode;
};

const useStyles = makeStyles((theme: Theme) => ({
	container: {
		width: "100%"
	},
	children: {
		maxWidth: "600px",
		margin: "0 auto",
		height: "100vh",
		display: "flex",
		justifyContent: "center",
		alignItems: "center",
		overflow: "scroll",
		padding: "0 1rem",
		scrollbarWidth: "none",
		msOverflowStyle: "none",
		"&::-webkit-scrollbar": {
			display: "none"
		},
		[`${theme.breakpoints.down("md")} and (orientation: landscape)`]: {
			height: "initial"
		}
	}
}));

const Layout: React.FC<LayoutProps> = ({ children }) => {
	const classes = useStyles();
	return (
		<div className={classes.container}>
			<div className={classes.children}>{children}</div>
		</div>
	);
};

export default Layout;
