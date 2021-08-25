/* eslint-disable @typescript-eslint/no-explicit-any */
import React from "react";
import { makeStyles, Theme } from "@material-ui/core/styles";
import Tabs from "@material-ui/core/Tabs";
import Tab from "@material-ui/core/Tab";
import TabPanel from "./TabPanel";
import CreateSecretForm from "./CreateSecretForm";
import { FormikHelpers, FormikValues } from "formik";
import { Box, Divider, Typography } from "@material-ui/core";

type SecretType = "file" | "message";

function a11yProps(index: any) {
	return {
		id: `form-tab-${index}`,
		"aria-controls": `simple-tabpanel-${index}`
	};
}

const useStyles = makeStyles((theme: Theme) => ({
	root: {
		maxWidth: "500px",
		margin: "0 auto",
		backgroundColor: theme.palette.background.paper,
		gap: theme.spacing(4),
		flexDirection: "column",
		justifyContent: "center",
		alignItems: "center",
		width: "100%",
		maxHeight: "80%",
		height: "80%",
		boxShadow: "rgba(0, 0, 0, 0.1) 0px 4px 12px",
		background: "#f5f7f8",
		borderRadius: "5px",
		padding: theme.spacing(2)
	},
	container: {}
}));

type CreateSecretFormTabsProps = {
	// eslint-disable-next-line @typescript-eslint/no-explicit-any
	onSubmit: (values: FormikValues, helpers: FormikHelpers<any>) => void;
	initialValues: FormikValues;
	loading: boolean;
};

const CreateSecretFormTabs: React.FC<CreateSecretFormTabsProps> = props => {
	const classes = useStyles();
	const [value, setValue] = React.useState<"file" | "message">("message");

	const handleChange = (event: React.ChangeEvent<any>, newValue: SecretType) => {
		setValue(newValue);
	};

	return (
		<div className={classes.root}>
			<Box textAlign="center" marginY="1rem">
				<Typography variant="body1" style={{ fontWeight: "bold" }}>
					Which kind of secret do you want to share ?
				</Typography>
			</Box>
			<Divider style={{ margin: "0 3rem" }} />
			<Tabs value={value} onChange={handleChange} centered>
				<Tab value="message" label="Secret message" {...a11yProps("message")} />
				<Tab value="file" label="Secret file" {...a11yProps("file")} />
			</Tabs>
			<TabPanel value={value} index="message">
				<CreateSecretForm
					type={value}
					onSubmit={props.onSubmit}
					initialValues={props.initialValues}
					loading={props.loading}
				/>
			</TabPanel>
			<TabPanel value={value} index="file">
				<CreateSecretForm
					type={value}
					onSubmit={props.onSubmit}
					initialValues={props.initialValues}
					loading={props.loading}
				/>
			</TabPanel>
		</div>
	);
};

export default CreateSecretFormTabs;
