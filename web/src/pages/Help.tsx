import { makeStyles } from "@material-ui/styles";
import Layout from "components/Layout";

const useStyles = makeStyles(() => ({
	container: {}
}));

const Help: React.FC = () => {
	const classes = useStyles();

	return (
		<Layout>
			<div className={classes.container}>help</div>
		</Layout>
	);
};

export default Help;
