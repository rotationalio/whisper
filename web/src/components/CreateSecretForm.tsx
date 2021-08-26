import { Typography } from "@material-ui/core";
import { TextField, Switch } from "formik-material-ui";
import { Field, Form, Formik } from "formik";
import { Autocomplete, AutocompleteRenderInputParams } from "formik-material-ui-lab";
import { TextField as MuiTextField } from "@material-ui/core";
import { Lifetime } from "utils/interfaces";
import React, { useEffect } from "react";
import { CreateSecretFormProps } from "types/CreateSecretFormProps";
import { preventNonNumericalInput } from "utils/utils";
import { useStyles } from "styles/createSecretFormStyles";
import clsx from "clsx";
import Button from "./Button";
import StyledDropzone from "./Dropzone";
import { CreateSecretFileSchema, CreateSecretMessageSchema } from "utils/validation-schema";
import { LIFETIME_OPTIONS } from "constants/index";

const CreateSecretForm: React.FC<CreateSecretFormProps> = props => {
	const classes = useStyles();
	const CreateSecretSchema = props.type === "file" ? CreateSecretFileSchema : CreateSecretMessageSchema;

	return (
		<Formik initialValues={props.initialValues} validationSchema={CreateSecretSchema} onSubmit={props.onSubmit}>
			{({ errors, values, setFieldValue }) => {
				useEffect(() => {
					if (values.accessType) {
						setFieldValue("accesses", -1);
					} else {
						setFieldValue("accesses", 1);
					}
				}, [values.accessType]);

				return (
					<Form className={classes.form} noValidate>
						<div className={clsx({ [classes.hide]: props.type === "file" })}>
							<Field
								component={TextField}
								name="secret"
								type="text"
								label="Type your secret here"
								variant="outlined"
								placeholder="Type your secret here"
								data-testid="secret"
								required
								multiline
								minRows={10}
								maxRows={15}
							/>
						</div>
						<div className={clsx({ [classes.hide]: props.type === "message" })}>
							<Field component={StyledDropzone} />
						</div>

						<div>
							<Field
								component={TextField}
								type="password"
								label="Your secret password"
								name="password"
								size="small"
								variant="outlined"
								helperText="Optional, set a password to unlock the secret"
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
								helperText="Limit the number of times the secret can be viewed"
								pattern="[0-9]*"
								InputProps={{ inputProps: { min: 1, max: 108 } }}
								onKeyPress={preventNonNumericalInput}
								onInput={(e: React.ChangeEvent<HTMLInputElement>) => {
									e.target.value = Math.max(-1, parseInt(e.target.value)).toString().slice(0, 3);
								}}
							/>
						</div>
						<div>
							<Field
								name="lifetime"
								component={Autocomplete}
								size="small"
								options={LIFETIME_OPTIONS}
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
							<Button
								label="Get Token"
								type="submit"
								variant="contained"
								color="primary"
								fullWidth
								isLoading={props.loading}
								disabled={props.loading}
							/>
						</div>
					</Form>
				);
			}}
		</Formik>
	);
};

export default CreateSecretForm;
