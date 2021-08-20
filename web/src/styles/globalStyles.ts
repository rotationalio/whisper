import { createStyles, makeStyles } from "@material-ui/core";

const useStyles = makeStyles(() =>
	createStyles({
		"@global": {}
	})
);

const GlobalStyles: React.FC = () => {
	useStyles();
	return null;
};

export default GlobalStyles;
