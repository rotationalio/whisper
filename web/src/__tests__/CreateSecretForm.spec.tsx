import { act, fireEvent, render, screen, waitFor } from "@testing-library/react";
import CreateSecretForm from "components/CreateSecretForm";
import userEvent from "@testing-library/user-event";

describe("CreateSecretForm", () => {
	describe("Should generate token with valid inputs", () => {
		it("should call submit function", async () => {
			const onSubmitMock = jest.fn();
			const initialValues = {
				secret: "",
				password: "",
				accessType: true,
				accesses: 1,
				lifetime: { value: "168h", label: "7 days" }
			};

			render(<CreateSecretForm onSubmit={onSubmitMock} initialValues={initialValues} />);

			const messageInput = screen.getByPlaceholderText(/type your secret here/i);
			const submitButton = screen.getByRole("button", { name: /get token/i });

			userEvent.type(messageInput, "the eagle flies at midnight");
			userEvent.click(submitButton);
			waitFor(() => expect(onSubmitMock).toHaveBeenCalledTimes(1));
		});
	});

	describe("Should not generate token with invalid inputs", () => {
		it("render the message validation error", async () => {
			const onSubmitMock = jest.fn();
			const { getByPlaceholderText, container } = render(
				<CreateSecretForm loading={false} type="message" onSubmit={onSubmitMock} initialValues={{}} />
			);
			const messageInput = getByPlaceholderText(/type your secret here/i);

			act(() => {
				fireEvent.blur(messageInput);
			});

			waitFor(() => expect(container.innerHTML).toMatch(/you must add a secret/i));
		});
	});
});
