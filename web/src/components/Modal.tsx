import React from "react";
import { IconButton, Modal as MuiModal } from "@material-ui/core";
import { useModal } from "contexts";
import CloseIcon from "@material-ui/icons/Close";
import dayjs from "dayjs";
import relativeTime from "dayjs/plugin/relativeTime";
import { useStyles } from "styles/modalStyles";
import { ModalType } from "utils/enums/modal";
dayjs.extend(relativeTime);

interface CustomModalProps {
	children: React.ReactNode;
}

const Modal: React.FC<CustomModalProps> = ({ children }) => {
	const classes = useStyles();
	const { state, dispatch } = useModal();

	const handleClose = () => {
		if (window.confirm("Be sure to copy the link before closing the window")) {
			dispatch({ type: ModalType.HIDE_MODAL });
		}
	};

	return (
		<MuiModal open={state.modalType === "SHOW_MODAL"} onClose={handleClose} aria-labelledby="token-modal">
			<div className={classes.paper}>
				<div>
					<IconButton className={classes.close} onClick={handleClose} title="Close">
						<CloseIcon />
					</IconButton>
				</div>
				{children}
			</div>
		</MuiModal>
	);
};

export default Modal;
