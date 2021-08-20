import { makeStyles, TextField, Theme, Typography } from "@material-ui/core";
import { Alert } from "@material-ui/lab";
import { FormikHelpers, FormikValues, useFormik } from "formik";
import * as Yup from "yup";
import React from "react";
import Button from "./Button";

const useStyles = makeStyles((theme: Theme) => ({
	container: {
		position: "absolute",
		top: "50%",
		left: "50%",
		transform: "translate(-50%, -50%)"
	},
	textfield: {
		width: "100%"
	},
	button__block: {
		marginTop: theme.spacing(3)
	}
}));

type SecretPasswordInput = {
	onSubmit: (values: FormikValues, helpers: FormikHelpers<{ password: string }>) => void;
	error: string;
};

const SecretPassword: React.FC<SecretPasswordInput> = ({ onSubmit, error }) => {
	const classes = useStyles();
	const formik = useFormik({
		initialValues: {
			password: ""
		},
		validationSchema: Yup.object().shape({
			password: Yup.string().required("A password is required to access this Secret")
		}),
		onSubmit: (values: FormikValues, helpers: FormikHelpers<{ password: string }>) => {
			onSubmit(values, helpers);
		}
	});

	return (
		<div className={classes.container}>
			<Alert severity="error" style={{ display: error ? undefined : "none", marginBottom: "1rem" }}>
				Password not accepted, please try again.
			</Alert>
			<Typography variant="body1" gutterBottom>
				This Secret requires a password for access.
			</Typography>
			<form onSubmit={formik.handleSubmit}>
				<div>
					<TextField
						type="password"
						label="Password"
						name="password"
						size="small"
						variant="outlined"
						onChange={formik.handleChange}
						value={formik.values.password}
						error={!!formik.errors["password"]}
						inputProps={{ "data-testid": "password-input" }}
						helperText={formik.errors["password"]}
						FormHelperTextProps={{
							title: "password error"
						}}
						className={classes.textfield}
					/>
				</div>
				<div className={classes.button__block}>
					<Button
						variant="contained"
						color="primary"
						type="submit"
						isLoading={formik.isSubmitting}
						label="Unlock Secret"
						style={{ maxWidth: "12rem", width: "100%" }}
					/>
				</div>
			</form>
		</div>
	);
};

export default SecretPassword;
