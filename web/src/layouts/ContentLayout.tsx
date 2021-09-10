import React from "react";
import { Chip, makeStyles, Theme } from "@material-ui/core";
import HelpOutlineIcon from "@material-ui/icons/HelpOutline";
import { useModal } from "contexts";
import { ModalType } from "utils/enums/modal";

const useStyles = makeStyles((theme: Theme) => ({
	root: {
		flexGrow: 1
	},
	container: {
		boxShadow: "rgba(0, 0, 0, 0.1) 0px 4px 12px",
		background: "#fff",
		padding: theme.spacing(2),
		borderRadius: "5px"
	},
	chpisContainer: {
		marginBottom: 5,
		textAlign: "end"
	}
}));

type ContentLayoutProps = {
	children: React.ReactNode;
};

const ContentLayout: React.FC<ContentLayoutProps> = ({ children }) => {
	const classes = useStyles();
	const { dispatch } = useModal();

	const handleClick = () => {
		dispatch({ type: ModalType.SHOW_ABOUT_US_MODAL });
	};

	return (
		<div className={classes.root}>
			<div className={classes.chpisContainer}>
				<Chip variant="outlined" size="small" avatar={<HelpOutlineIcon />} label="Learn More" onClick={handleClick} />
			</div>
			<div className={classes.container}>{children}</div>
		</div>
	);
};

export default ContentLayout;
