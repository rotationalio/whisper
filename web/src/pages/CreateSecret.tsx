import React from "react";
import { Grid } from "@material-ui/core";
import createSecret from "services/createSecret";
import { Secret } from "utils/interfaces/Secret";
import { FormikValues } from "formik";
import { AxiosError, AxiosResponse } from "axios";
import { useStyles } from "styles/createSecretStyles";
import { Alert, Color } from "@material-ui/lab";
import { Lifetime } from "utils/interfaces";
import { useModal } from "contexts/modalContext";
import { ModalType } from "utils/enums/modal";
import Layout from "components/Layout";
import { encodeFileToBase64 } from "utils/utils";
import CreateSecretFormTabs from "components/CreateSecretFormTabs";

interface Values {
	secret: string;
	password: string;
	accessType: boolean;
	accesses: number;
	lifetime: Lifetime;
	file: null;
}

const initialValues: Values = {
	secret: "",
	password: "",
	accessType: true,
	accesses: 1,
	file: null,
	lifetime: { value: "168h", label: "7 days" }
};

const CreateSecret: React.FC = () => {
	const [, setToken] = React.useState<{ token: string; expires: Date }>();
	const [message, setMessage] = React.useState<{ status?: Color; message?: string }>({
		status: undefined,
		message: undefined
	});
	const [isLoading, setIsLoading] = React.useState(false);
	const { dispatch } = useModal();

	const classes = useStyles();

	async function handleSubmit(values: FormikValues) {
		setIsLoading(true);
		const lifetime = values.lifetime ? values.lifetime.value : { value: "168h", label: "7 days" };
		const encodedSecretFile = values.file ? await encodeFileToBase64(values.file) : "";

		const data: Secret = {
			lifetime,
			secret: values.secret || encodedSecretFile,
			password: values.password,
			accesses: values.accesses,
			filename: values.filename || "",
			is_base64: values.is_base64 || false
		};

		createSecret(data).then(
			(response: AxiosResponse) => {
				setToken(response.data);
				setIsLoading(false);

				dispatch({ type: ModalType.SHOW_MODAL, payload: response.data });
			},
			(error: AxiosError) => {
				setMessage({ status: "error", message: error.message });
				setTimeout(() => {
					setMessage({ status: undefined, message: undefined });
				}, 5000);

				setIsLoading(false);
			}
		);
	}

	return (
		<Layout>
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
						<CreateSecretFormTabs onSubmit={handleSubmit} initialValues={initialValues} loading={isLoading} />
					</Grid>
				</Grid>
			</div>
		</Layout>
	);
};

export default CreateSecret;
