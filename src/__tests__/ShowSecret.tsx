import ShowSecret from "components/ShowSecret";
import dayjs from "dayjs";
import relativeTime from "dayjs/plugin/relativeTime";
import { render } from "utils/test-utils";
dayjs.extend(relativeTime);

describe("ShowSecret", () => {
	const token = "iwbNM2NqW93FcKzl1FBVR0awIC41nemQhXdgE4ul-PA";
	const secretMock = {
		secret: "the eagle flies at midnight",
		is_base64: false,
		created: new Date("2021-08-16T16:48:14.584642555Z"),
		accesses: 1,
		destroyed: false,
		lifetime: "10m"
	};

	it("should render secret message", () => {
		const { container } = render(<ShowSecret secret={secretMock} token={token} />);
		expect(container).toHaveTextContent("the eagle flies at midnight");
	});

	it("should render secret creation date", () => {
		const { container } = render(<ShowSecret secret={secretMock} token={token} />);
		const formatedDate = dayjs(secretMock.created).fromNow();
		expect(container).toHaveTextContent(`${formatedDate}`);
	});
});
