import FetchSecret from "pages/FetchSecret";
import { renderWithRouterMatch, waitFor } from "utils/test-utils";
import { SecretMock } from "__mocks__/SecretMock";
import { getSecret } from "../services";

jest.mock("../services");

const mockedGetSecret = getSecret as jest.MockedFunction<typeof getSecret>;

afterEach(() => {
	jest.clearAllMocks();
});
describe("FetchSecret", () => {
	it("should display ShowSecret component", async () => {
		mockedGetSecret.mockResolvedValueOnce(SecretMock as any);

		renderWithRouterMatch(<FetchSecret />, {
			path: "/secret/:token",
			route: "/secret/iwbNM2NqW93FcKzl1FBVR0awIC41nemQhXdgE4ul-PA"
		});

		await waitFor(() => {
			expect(mockedGetSecret).toHaveBeenCalledTimes(1);
			expect(mockedGetSecret).toHaveBeenCalledWith("iwbNM2NqW93FcKzl1FBVR0awIC41nemQhXdgE4ul-PA");
		});
	});
});
