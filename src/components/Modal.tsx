import React, { useState } from "react";
import { Box, Button, CircularProgress, IconButton, Modal as MuiModal, Snackbar, Typography } from "@material-ui/core";
import { useModal } from "contexts";
import LinkIcon from "@material-ui/icons/Link";
import CloseIcon from "@material-ui/icons/Close";
import { CopyToClipboard } from "react-copy-to-clipboard";
import dayjs from "dayjs";
import relativeTime from "dayjs/plugin/relativeTime";
import { Alert, Color } from "@material-ui/lab";
import deleteSecret from "services/deleteSecret";
import { AxiosError } from "axios";
import { useStyles } from "styles/modalStyles";
import { ModalType } from "utils/enums/modal";
import { generateSecretLink } from "utils/utils";
dayjs.extend(relativeTime);

const Modal: React.FC = () => {
	const classes = useStyles();
	const { state, dispatch } = useModal();
	const [alert, setAlert] = useState<{ open: boolean; severity?: Color; message?: string }>({
		open: false,
		severity: "success",
		message: ""
	});
	const [isLoading, setIsLoading] = useState(false);
	const secretLink = typeof state.modalProps?.token === "string" ? generateSecretLink(state.modalProps?.token) : "";

	const handleCopy = () => setAlert({ open: true, message: "Secret link copied" });
	const handleAlertClose = () => setAlert({ open: false });

	const handleClose = () => {
		if (window.confirm("Be sure to copy the link before closing the window")) {
			dispatch({ type: ModalType.HIDE_MODAL });
		}
	};

	const handleDeleteSecret = async () => {
		if (window.confirm("Do you really want to remove the secret message ?")) {
			if (state.modalProps?.token) {
				setIsLoading(true);
				deleteSecret(state.modalProps?.token).then(
					() => {
						setIsLoading(false);
						setAlert({ open: true, message: "Secret message destroyed" });
						dispatch({ type: ModalType.HIDE_MODAL });
					},
					(error: AxiosError) => {
						setIsLoading(false);
						setAlert({ open: true, message: error.message, severity: "error" });
					}
				);
			}
		}
	};

	return (
		<div>
			<MuiModal open={state.modalType === "SHOW_MODAL"} onClose={handleClose} aria-labelledby="token-modal">
				<div className={classes.paper}>
					<div>
						<IconButton className={classes.close} onClick={handleClose} title="Close">
							<CloseIcon />
						</IconButton>
					</div>
					<Box>
						<Typography variant="h5" align="center" gutterBottom>
							Secret created successfully
						</Typography>
						<Typography align="center" gutterBottom>
							You can find your secret on this link below
						</Typography>
						<Typography align="center" gutterBottom>
							It expires{" "}
							<span style={{ color: "red", fontWeight: "bold" }}>{dayjs(state.modalProps?.expires).fromNow()}</span>
						</Typography>
					</Box>
					<Box
						display="flex"
						alignItems="center"
						paddingLeft="15px"
						justifyContent="space-between"
						bgcolor="#f1f3f4"
						borderRadius=".2rem"
					>
						<div className={classes.ellipsis}>
							<Typography>{secretLink}</Typography>
						</div>
						<div>
							<CopyToClipboard text={secretLink} onCopy={handleCopy}>
								<IconButton aria-label="copy" size="medium" title="Copy this link">
									<LinkIcon />
								</IconButton>
							</CopyToClipboard>
						</div>
					</Box>
					<Box display="flex" gridGap="2rem" justifyContent="space-around">
						<Button
							onClick={handleDeleteSecret}
							color="secondary"
							variant="contained"
							disabled={isLoading}
							style={{ minWidth: "200px" }}
						>
							{isLoading ? <CircularProgress size={24} /> : "Destroy this secret"}
						</Button>
					</Box>
				</div>
			</MuiModal>
			<Snackbar open={alert.open} autoHideDuration={5000} onClose={handleAlertClose}>
				<Alert onClose={handleAlertClose} severity={alert.severity}>
					{alert.message}
				</Alert>
			</Snackbar>
		</div>
	);
};

export default Modal;
