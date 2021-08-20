import { AxiosError, AxiosResponse } from "axios";
import SecretNotFound from "components/SecretNotFound";
import { FormikHelpers, FormikValues } from "formik";
import React from "react";
import { useParams } from "react-router-dom";
import getSecret from "services/ShowSecret";
import { Secret } from "utils/interfaces/Secret";
import SecretPassword from "../components/SecretPassword";
import dayjs from "dayjs";
import relativeTime from "dayjs/plugin/relativeTime";
import ShowSecret from "components/ShowSecret";
import Layout from "components/Layout";
dayjs.extend(relativeTime);

const FetchSecret: React.FC = () => {
	const [secret, setSecret] = React.useState<Secret>();
	const [status, setStatus] = React.useState("pending");
	const [errorMessage, setErrorMessage] = React.useState("");
	const isMounted = React.useRef(true);
	const { token } = useParams<{ token: string }>();

	React.useEffect(() => {
		if (isMounted) {
			getSecret(token).then(
				(response: AxiosResponse) => {
					setSecret(response.data);
					setStatus("success");
				},
				(error: AxiosError) => {
					if (error.response?.status === 401) {
						setStatus("unauthorized");
					} else if (error.response?.status === 404) {
						setStatus("error");
						setErrorMessage("No secret exists with the specified token.");
					}
				}
			);
		}
		return () => {
			isMounted.current = false;
		};
	}, []);

	const handlePasswordSubmit = async (values: FormikValues, helpers: FormikHelpers<{ password: string }>) => {
		getSecret(token, values.password).then(
			(response: AxiosResponse) => {
				setSecret(response.data);
				setStatus("success");
				helpers.setSubmitting(false);
			},
			(error: AxiosError) => {
				if (error.response && error.response?.status !== 401) {
					setStatus("error");
					setErrorMessage("No secret exists with the specified token.");
				}
				setErrorMessage(error.message);
				helpers.setSubmitting(false);
			}
		);
	};

	return (
		<Layout>
			{status === "pending" ? "Loading..." : null}
			{status === "error" ? <SecretNotFound title="Not Found" message={errorMessage} /> : null}
			{status === "unauthorized" ? <SecretPassword onSubmit={handlePasswordSubmit} error={errorMessage} /> : null}
			{status === "success" ? <ShowSecret secret={secret} token={token} /> : null}
		</Layout>
	);
};

export default FetchSecret;
