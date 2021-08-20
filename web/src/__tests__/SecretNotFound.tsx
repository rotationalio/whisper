import SecretNotFound from "components/SecretNotFound";
import { render, screen } from "utils/test-utils";

describe("SecretNotFound", () => {
	it("Should render not found message", () => {
		render(<SecretNotFound message="Secret does not exist in secret manager." title="Not found" />);
		expect(screen.getByRole("alert")).toBeInTheDocument();

		expect(screen.getByText(/Secret does not exist in secret manager./i)).toBeInTheDocument();
		expect(screen.getByText(/Not found/i)).toBeInTheDocument();
	});
});
