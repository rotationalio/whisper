import { Button, makeStyles, Theme, Typography } from "@material-ui/core";
import { TextField, Switch } from "formik-material-ui";
import { Field, Form, Formik } from "formik";
import { Autocomplete, AutocompleteRenderInputParams } from "formik-material-ui-lab";
import { TextField as MuiTextField } from "@material-ui/core";
import * as Yup from "yup";
import { Lifetime } from "utils/interfaces";
import { useEffect } from "react";

const useStyles = makeStyles((theme: Theme) => ({
	form: {
		display: "flex",
		flexDirection: "column",
		gap: theme.spacing(2),
		maxWidth: 400,
		margin: "auto",
		"& .MuiTextField-root": {
			width: "100%"
		}
	},
	height__full: {
		height: "100%"
	},
	hide: {
		display: "none"
	}
}));

const options: Lifetime[] = [
	{ value: "5m", label: "5 min" },
	{ value: "15m", label: "15 min" },
	{ value: "30m", label: "30 min" },
	{ value: "1h", label: "1 hour" },
	{ value: "2h", label: "2 hours" },
	{ value: "3h", label: "3 hours" },
	{ value: "24h", label: "1 day" },
	{ value: "48h", label: "2 days" },
	{ value: "72h", label: "3 days" },
	{ value: "168h", label: "7 days" }
];

const initialValues = {
	secret: "",
	password: "",
	accessType: true,
	accessNumber: 1,
	lifetime: { value: "7d", label: "7 days" }
};

const CreateSecretSchema = Yup.object().shape({
	secret: Yup.string().required("You must add a secret"),
	password: Yup.string(),
	lifetime: Yup.object().shape({ value: Yup.string(), label: Yup.string() }).nullable(),
	accessType: Yup.boolean(),
	accessNumber: Yup.number().max(108, "The max number is 108")
});

const CreateSecretForm: React.FC = () => {
	const classes = useStyles();

	function preventNonNumericalInput(e: any) {
		e = e || window.event;
		const charCode = typeof e.which == "undefined" ? e.keyCode : e.which;
		const charStr = String.fromCharCode(charCode);

		if (!charStr.match(/^[0-9]+$/)) e.preventDefault();
	}

	return (
		<Formik
			initialValues={initialValues}
			validationSchema={CreateSecretSchema}
			onSubmit={(values, { setSubmitting }) => {
				setTimeout(() => {
					setSubmitting(false);
					alert(JSON.stringify(values, null, 2));
				}, 200);
			}}
		>
			{({ isSubmitting, errors, values, setFieldValue }) => {
				useEffect(() => {
					if (values.accessType) {
						setFieldValue("accessNumber", -1);
					} else {
						setFieldValue("accessNumber", 1);
					}
				}, [values.accessType]);

				return (
					<Form className={classes.form} noValidate>
						<div>
							<Field
								component={TextField}
								name="secret"
								type="text"
								label="Type your secret here"
								variant="outlined"
								placeholder="Type your secret here"
								required
								multiline
								minRows={10}
								maxRows={15}
							/>
						</div>
						<div>
							<Field
								component={TextField}
								type="password"
								label="Your secret password"
								name="password"
								size="small"
								variant="outlined"
								disabled
							/>
						</div>
						<div>
							<Field component={Switch} type="checkbox" name="accessType" />
							<Typography variant="button">unlimited access</Typography>
						</div>
						<div className={values["accessType"] ? classes.hide : undefined}>
							<Field
								component={TextField}
								type="number"
								label="Number of access"
								name="accessNumber"
								size="small"
								variant="outlined"
								pattern="[0-9]*"
								InputProps={{ inputProps: { min: 1, max: 108 } }}
								onKeyPress={preventNonNumericalInput}
								onInput={(e: any) => {
									e.target.value = Math.max(-1, parseInt(e.target.value)).toString().slice(0, 3);
								}}
							/>
						</div>
						<div>
							<Field
								name="lifetime"
								component={Autocomplete}
								size="small"
								options={options}
								disabled
								getOptionLabel={(option: Lifetime) => option.label}
								getOptionSelected={(option: Lifetime, value: Lifetime) => option.value === value.value}
								renderInput={(params: AutocompleteRenderInputParams) => (
									<MuiTextField
										{...params}
										label="Secret lifetime"
										variant="outlined"
										error={!!errors["lifetime"]}
										helperText={
											errors["lifetime"]
												? errors["lifetime"]
												: "By default, the secret has 7 days before being destroyed"
										}
									/>
								)}
							/>
						</div>
						<div>
							<Button type="submit" variant="contained" color="primary" disabled={isSubmitting}>
								Get Token
							</Button>
						</div>
					</Form>
				);
			}}
		</Formik>
	);
};

export default CreateSecretForm;
