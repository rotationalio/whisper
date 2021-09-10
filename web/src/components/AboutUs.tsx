import { Box, Typography } from "@material-ui/core";
import { useModal } from "contexts";
import { useServerStatus } from "contexts/serverStatusContext";
import React from "react";
import { ModalType } from "utils/enums/modal";
import Badge from "./Badge";
import Modal from "./Modal";

const AboutUs: React.FC = () => {
	const { state, dispatch } = useModal();
	const [status] = useServerStatus();

	const StatusColor = {
		ok: "green",
		maintainance: "yellow",
		unhealthy: "red"
	};

	const handleClose = () => dispatch({ type: ModalType.HIDE_MODAL });

	return (
		<Modal open={state.modalType === ModalType.SHOW_ABOUT_US_MODAL} onClose={handleClose}>
			<>
				<div>
					<h3>About Whisper</h3>
					<p>
						The Whisper service is an internal helper tool used at <a href="https://rotational.io">Rotational Labs</a>{" "}
						to quickly share secrets, configurations, environment files, credentials, certificates, and more. Whisper is
						designed to accelerate our own internal software engineering practice and is comprised of an API service
						that is accessed by both a web UI and a command line application. There are many tools like Whisper, but
						this one is ours!
					</p>
					<p>
						To download the CLI application, report bugs or issues, or learn more about Whisper, please see the
						README.md file in the Whisper GitHub repository:{" "}
						<a href="https://github.com/rotationalio/whisper">rotationalio/whisper</a>.
					</p>
					<p>
						Although Whisper is an internal tool at Rotational, We&apos;ve made the code open source and are happy to
						have general contributions that enhance the project (particularly if you&apos;re a member of the Rotational
						Engineering Team!) We&apos;ve made our releases and the code freely available under the{" "}
						<a href="https://github.com/rotationalio/whisper/blob/main/LICENSE">Apache License 2.0</a> and we&apos;d
						feel privileged if you used Whisper in your own organization. Please note, however, that Rotational Labs
						makes no guarantees about the security of this software project and provides all code and binaries as is for
						general use. Use with common sense and at your own risk!
					</p>
					<p>
						If you&apos;re a Rotational customer and are interested in Whisper, please let us know, we&apos;d be happy
						to deploy it for you as a single-tenant service. If you&apos;re not a Rotational customer but are
						interested, please get in touch with us at <a href="mailto:info@rotational.io">info@rotational.io</a>.
					</p>
				</div>
				<footer>
					<Box display="flex" flexDirection="column" alignItems="center">
						<Box display="flex" alignItems="center" gridGap=".5rem" marginBottom=".3rem">
							{status.status && status.version && (
								<>
									<Typography variant="caption">
										{status.status && (
											<span>
												<Badge content={`${status?.host || "server"} status`} color={StatusColor[status?.status]} />
											</span>
										)}
									</Typography>
									<Typography variant="caption">
										version: <span style={{ fontWeight: "bold" }}>{status?.version}</span>
									</Typography>
								</>
							)}
						</Box>
						<Typography>
							Made with &spades; by <a href="https://rotational.io/">Rotational Labs</a>
						</Typography>
					</Box>
				</footer>
			</>
		</Modal>
	);
};

export default AboutUs;
