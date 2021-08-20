import userEvent from "@testing-library/user-event";
import SecretPassword from "components/SecretPassword";
import { fireEvent, render, renderWithRouterMatch, screen, waitFor } from "utils/test-utils";

describe("SecretPassword", () => {
	const handleSubmit = jest.fn();

	it("should show validation on submit", async () => {
		render(<SecretPassword onSubmit={handleSubmit} error="wrong password" />);

		const submitButton = screen.getByRole("button", { name: /unlock secret/i });
		fireEvent.click(submitButton);

		await waitFor(() => {
			expect(screen.getByTitle(/password error/i)).toBeInTheDocument();
			expect(screen.getByTitle(/password error/i)).toHaveClass("Mui-error");
		});
	});

	it("Should call handleSubmit", async () => {
		renderWithRouterMatch(<SecretPassword onSubmit={handleSubmit} error="wrong password" />, {
			path: "/secret/:token",
			route: "/secret/_0lzPfLGsL11wWKRRefVvQb5hIn12Ln2n4PQc4H25Fs"
		});
		const passwordInput = screen.getByTestId(/password-input/) as HTMLInputElement;
		const submitButton = screen.getByRole("button", { name: /unlock secret/i });

		fireEvent.change(passwordInput, { target: { value: "password" } });
		userEvent.click(submitButton);

		await waitFor(() => {
			expect(handleSubmit).toHaveBeenCalledTimes(1);
		});
	});
});
