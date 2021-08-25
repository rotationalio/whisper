import { FormikHelpers, FormikValues } from "formik";

export type CreateSecretFormProps = {
	// eslint-disable-next-line @typescript-eslint/no-explicit-any
	onSubmit: (values: FormikValues, helpers: FormikHelpers<any>) => void;
	initialValues: FormikValues;
	loading: boolean;
	type: "file" | "message";
};
