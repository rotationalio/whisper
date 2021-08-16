import { FormikHelpers, FormikValues } from "formik";

export type CreateSecretFormProps = {
	onSubmit: (values: FormikValues, helpers: FormikHelpers<any>) => void;
	initialValues: FormikValues;
};
