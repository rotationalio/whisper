import React from "react";
import { Grid, makeStyles, Theme } from "@material-ui/core";
import CreateSecretForm from "components/CreateSecretForm";

const useStyles = makeStyles((theme: Theme) => ({
	root: {
		flexGrow: 1,
		height: "100vh",
		padding: theme.spacing(2)
	},
	h__full: {
		height: "100%"
	},
	w__full: {
		width: "100%"
	}
}));

const CreateSecret: React.FC = () => {
	const classes = useStyles();

	return (
		<div className={classes.root}>
			<Grid container alignItems="center" className={classes.h__full}>
				<Grid item sm={12} md={12} className={classes.w__full}>
					<CreateSecretForm />
				</Grid>
			</Grid>
		</div>
	);
};

export default CreateSecret;
