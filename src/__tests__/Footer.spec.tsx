import { render } from "@testing-library/react";
import Footer from "components/Footer";

describe("Footer", () => {
	it("should render correctly", () => {
		const { asFragment } = render(<Footer />);

		expect(asFragment()).toMatchInlineSnapshot(`
		<DocumentFragment>
		  <footer
		    class="makeStyles-root-1"
		  >
		    <p
		      class="MuiTypography-root MuiTypography-body1"
		    >
		      Made with ♠ by 
		      <a
		        class="MuiTypography-root MuiLink-root MuiLink-underlineHover makeStyles-text__white-2 MuiTypography-colorPrimary"
		        href="https://rotational.io"
		        target="_blank"
		      >
		        Rotational Labs
		      </a>
		    </p>
		    <div
		      class="MuiBox-root MuiBox-root-3"
		    >
		      <div
		        aria-label="add"
		        class="MuiBox-root MuiBox-root-4"
		        style="cursor: pointer;"
		        title="connected to "
		      >
		        <span
		          class="MuiTypography-root MuiTypography-caption"
		        >
		          status
		          <span
		            class="makeStyles-root-5 makeStyles-root-6"
		            color="green"
		          />
		        </span>
		        <span
		          class="MuiTypography-root MuiTypography-caption"
		        >
		          version: 
		        </span>
		      </div>
		    </div>
		  </footer>
		</DocumentFragment>
	`);
	});
});
