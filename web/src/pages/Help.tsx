import { makeStyles } from "@material-ui/styles";
import { GlobalLayout } from "layouts";

const useStyles = makeStyles(() => ({
	container: {}
}));

const Help: React.FC = () => {
	const classes = useStyles();

	return (
		<GlobalLayout>
			<div className={classes.container}>help</div>
		</GlobalLayout>
	);
};

export default Help;
