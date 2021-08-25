import { makeStyles } from "@material-ui/core";
import Footer from "./Footer";

type LayoutProps = {
	children: React.ReactNode;
};

const useStyles = makeStyles(() => ({
	children: {
		height: "100vh",
		width: "100%"
	}
}));

const Layout: React.FC<LayoutProps> = ({ children }) => {
	const classes = useStyles();
	return (
		<div>
			<div className={classes.children}>{children}</div>
			<Footer />
		</div>
	);
};

export default Layout;
