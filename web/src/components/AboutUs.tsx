import { Box, Tooltip, Typography } from "@material-ui/core";
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
					Lorem ipsum dolor sit amet consectetur adipisicing elit. Soluta voluptate iste maiores praesentium rerum eos
					pariatur suscipit animi alias similique, ad qui, minus officia aspernatur dicta cumque fugit? Ab vel
					necessitatibus dolor nesciunt pariatur inventore dolores corrupti beatae, laboriosam, nobis architecto
					sapiente iusto omnis aliquam asperiores officia cumque, doloribus non labore a. Deserunt nobis corporis
					deleniti harum delectus facere beatae vitae ducimus veritatis quas, esse doloribus maiores molestiae possimus
					a voluptatibus, velit dolores pariatur saepe debitis, enim error quasi similique! Debitis aut asperiores,
					aliquid harum veniam quia modi non tempora, illum sequi labore amet consectetur corrupti facere, similique sed
					aliquam in! Accusantium rem commodi autem sit fugit corporis id aperiam fuga facere obcaecati quo omnis, nisi
					reprehenderit quaerat et velit, molestias labore. Mollitia dolore doloremque explicabo architecto culpa,
					sapiente asperiores inventore molestiae aut incidunt fuga ipsa quis numquam illum molestias animi eius aliquam
					rerum. Sunt totam quis nostrum aspernatur inventore numquam eaque ducimus architecto eum officia maxime
					distinctio, dolorem atque magni quisquam ipsa consequuntur aliquid ad ratione rerum aut! Illum, possimus velit
					nemo, minima assumenda sapiente cum dolores sed ullam laboriosam cupiditate quo et eveniet quae! Molestiae
					error laudantium pariatur placeat culpa beatae numquam, asperiores at, saepe facilis optio sequi?
				</div>
				<footer>
					<Box display="flex" flexDirection="column" alignItems="center">
						<Box display="flex" alignItems="center" gridGap=".5rem" marginBottom=".3rem">
							{status.status && status.version && (
								<>
									<Typography variant="caption">
										{status.status && (
											<Tooltip title={status.host || ""} style={{ cursor: "pointer" }}>
												<span>
													<Badge content="server" color={StatusColor[status?.status]} />
												</span>
											</Tooltip>
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
