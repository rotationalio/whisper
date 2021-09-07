import React from "react";
import { IconButton, Modal as MuiModal, ModalProps } from "@material-ui/core";
import CloseIcon from "@material-ui/icons/Close";
import dayjs from "dayjs";
import relativeTime from "dayjs/plugin/relativeTime";
import { useStyles } from "styles/modalStyles";
dayjs.extend(relativeTime);

interface CustomModalProps extends ModalProps {
	children: React.ReactElement;
	onClose: () => void;
}

const Modal: React.FC<CustomModalProps> = ({ children, onClose, ...props }) => {
	const classes = useStyles();

	return (
		<MuiModal aria-labelledby="token-modal" {...props}>
			<div className={classes.paper}>
				<div>
					<IconButton className={classes.close} onClick={onClose} title="Close">
						<CloseIcon />
					</IconButton>
				</div>
				{children}
			</div>
		</MuiModal>
	);
};

export default Modal;
