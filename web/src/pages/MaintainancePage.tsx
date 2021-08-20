import { Box, makeStyles, Theme, Typography } from "@material-ui/core";

const useStyles = makeStyles((theme: Theme) => ({
	container: {
		position: "absolute",
		top: "50%",
		left: "50%",
		transform: "translate(-50%, -50%)",
		color: "#222121",
		padding: theme.spacing(2)
	},
	fontBold: {
		fontWeight: "bold"
	},
	textGray: {}
}));

const MaintainancePage: React.FC = () => {
	const classes = useStyles();

	return (
		<div className={classes.container}>
			<Box maxWidth="500px">
				<Typography variant="h3" className={classes.fontBold}>
					We&apos;ll be back soon !
				</Typography>
				<Typography>Sorry for the inconvenience but we&apos;re performing some maintainance at the moment</Typography>
			</Box>
		</div>
	);
};

export default MaintainancePage;
