import { render } from "@testing-library/react";
import Badge from "components/Badge";

describe("Badge", () => {
	it("should render badge content", () => {
		const { container } = render(<Badge color="red" content="badge content" />);
		expect(container.firstChild).toHaveTextContent("badge content");
	});

	it("should render correctly", () => {
		const { asFragment } = render(<Badge color="red" content="badge content" />);

		expect(asFragment()).toMatchInlineSnapshot(`
		<DocumentFragment>
		  <span
		    class="MuiTypography-root MuiTypography-caption"
		  >
		    badge content
		    <span
		      class="makeStyles-root-3 makeStyles-root-4"
		      color="red"
		    />
		  </span>
		</DocumentFragment>
	`);
	});
});
