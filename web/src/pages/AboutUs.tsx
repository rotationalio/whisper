import { makeStyles } from "@material-ui/styles";
import { GlobalLayout } from "layouts";

const useStyles = makeStyles(() => ({
	container: {}
}));

const AboutUs: React.FC = () => {
	const classes = useStyles();

	return (
		<GlobalLayout>
			<div className={classes.container}>about</div>
		</GlobalLayout>
	);
};

export default AboutUs;
