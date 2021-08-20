import { makeStyles, Typography } from "@material-ui/core";
import React from "react";
import { Link } from "react-router-dom";

const useStyle = makeStyles(() => ({
	container: {
		position: "absolute",
		top: "50%",
		left: "50%",
		transform: "translate(-50%, -50%)",
		textAlign: "center"
	}
}));

const NotFound: React.FC = () => {
	const classes = useStyle();
	return (
		<div className={classes.container}>
			<Typography variant="h1" style={{ fontWeight: "bold" }}>
				404
			</Typography>
			<Typography variant="h6" gutterBottom>
				This page could not be found
			</Typography>
			<Link to="/">Back to home</Link>
		</div>
	);
};

export default NotFound;
