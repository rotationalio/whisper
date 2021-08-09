import React, { useState } from "react";
import { Grid } from "@material-ui/core";
import CreateSecretForm from "components/CreateSecretForm";
import createSecret from "services/createSecret";
import { Secret } from "utils/interfaces/Secret";
import { FormikHelpers, FormikValues } from "formik";
import { AxiosError, AxiosResponse } from "axios";
import * as Yup from "yup";
import { useStyles } from "styles/createSecretStyles";
import { Alert, Color } from "@material-ui/lab";
import { Lifetime } from "utils/interfaces";
import { useModal } from "contexts/modalContext";
import { ModalType } from "utils/enums/modal";

const CreateSecretSchema = Yup.object().shape({
	secret: Yup.string().required("You must add a secret"),
	password: Yup.string(),
	lifetime: Yup.object().shape({ value: Yup.string(), label: Yup.string() }).nullable(),
	accessType: Yup.boolean(),
	accesses: Yup.number().max(108, "The max number is 108")
});

interface Values {
	secret: string;
	password: string;
	accessType: boolean;
	accesses: number;
	lifetime: Lifetime;
}

const initialValues: Values = {
	secret: "",
	password: "",
	accessType: true,
	accesses: 1,
	lifetime: { value: "168h", label: "7 days" }
};

const CreateSecret: React.FC = () => {
	const [, setToken] = useState<{ token: string; expires: Date }>();
	const [message, setMessage] = useState<{ status?: Color; message?: string }>({
		status: undefined,
		message: undefined
	});
	const { dispatch } = useModal();

	const classes = useStyles();

	function handleSubmit(values: FormikValues, helpers: FormikHelpers<Values>) {
		const data: Secret = {
			lifetime: values.lifetime.value,
			secret: values.secret,
			password: values.password,
			accesses: values.accesses,
			filename: values.filename || "",
			is_base64: false
		};

		createSecret(data).then(
			(response: AxiosResponse) => {
				setToken(response.data);
				helpers.setSubmitting(false);

				dispatch({ type: ModalType.SHOW_MODAL, payload: response.data });
			},
			(error: AxiosError) => {
				console.error("[error when creating secret]", error.message);
				setMessage({ status: "error", message: error.message });
				setTimeout(() => {
					setMessage({ status: undefined, message: undefined });
				}, 5000);

				helpers.setSubmitting(false);
			}
		);
	}

	return (
		<div className={classes.root}>
			<Grid container alignItems="center" className={classes.h__full}>
				<Grid item sm={12} md={12} className={classes.w__full}>
					<Alert
						severity={message.status}
						className={classes.alert}
						style={{ visibility: message.message ? "visible" : "hidden" }}
					>
						{message.message}
					</Alert>
					<CreateSecretForm
						onSubmit={handleSubmit}
						initialValues={initialValues}
						validationSchema={CreateSecretSchema}
					/>
				</Grid>
			</Grid>
		</div>
	);
};

export default CreateSecret;
