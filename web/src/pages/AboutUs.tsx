import { makeStyles } from "@material-ui/styles";
import Layout from "components/Layout";

const useStyles = makeStyles(() => ({
	container: {}
}));

const AboutUs: React.FC = () => {
	const classes = useStyles();

	return (
		<Layout>
			<div className={classes.container}>about</div>
		</Layout>
	);
};

export default AboutUs;
