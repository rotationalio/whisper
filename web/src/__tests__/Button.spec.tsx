import { render, fireEvent, screen } from "@testing-library/react";
import Button from "components/Button";

describe("Button", () => {
	it("should render button content", () => {
		const { container } = render(<Button label="Submit" />);
		expect(container).toHaveTextContent("Submit");
	});

	it("should render children", () => {
		const { container } = render(<Button label="Submit">Hello</Button>);
		expect(container).toHaveTextContent("Hello");
	});

	it("Should be clickable", () => {
		const handleClick = jest.fn();
		render(<Button label="Submit" onClick={handleClick} />);
		const button = screen.getByTestId("custom-button");
		fireEvent.click(button);

		expect(handleClick).toHaveBeenCalled();
	});
});
