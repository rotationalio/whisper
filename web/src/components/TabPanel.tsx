/* eslint-disable @typescript-eslint/no-explicit-any */
import { Box } from "@material-ui/core";

interface TabPanelProps {
	children?: React.ReactNode;
	index: any;
	value: any;
}

const TabPanel: React.FC<TabPanelProps> = props => {
	const { children, value, index, ...other } = props;

	return (
		<div
			role="createformpanel"
			hidden={value !== index}
			id={`form-tabpanel-${index}`}
			aria-labelledby={`form-tab-${index}`}
			{...other}
		>
			{value === index && <Box p={3}>{children}</Box>}
		</div>
	);
};

export default TabPanel;
