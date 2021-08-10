import { Button, Typography } from "@material-ui/core";
import { TextField, Switch } from "formik-material-ui";
import { Field, Form, Formik } from "formik";
import { Autocomplete, AutocompleteRenderInputParams } from "formik-material-ui-lab";
import { TextField as MuiTextField } from "@material-ui/core";
import { Lifetime } from "utils/interfaces";
import { useEffect } from "react";
import { CreateSecretFormProps } from "types/CreateSecretFormProps";
import { preventNonNumericalInput } from "utils/utils";
import { useStyles } from "styles/createSecretFormStyles";

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

const CreateSecretForm: React.FC<CreateSecretFormProps> = props => {
	const classes = useStyles();

	return (
		<Formik initialValues={props.initialValues} validationSchema={props.validationSchema} onSubmit={props.onSubmit}>
			{({ isSubmitting, errors, values, setFieldValue }) => {
				useEffect(() => {
					if (values.accessType) {
						setFieldValue("accesses", -1);
					} else {
						setFieldValue("accesses", 1);
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
								name="accesses"
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
